package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
)

type Summary struct {
    Tables  int
    Rows    int
    Incons  int
    Skipped int
    Files   int
}

func (s *Summary) Report() string {
    const (
        reset  = "\x1b[0m"
        green  = "\x1b[32m"
        yellow = "\x1b[33m"
        red    = "\x1b[31m"
    )
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%sTables processed:%s %d\n", green, reset, s.Tables))
    sb.WriteString(fmt.Sprintf("%sRows processed:%s %d\n", green, reset, s.Rows))
    sb.WriteString(fmt.Sprintf("%sInconsistent rows:%s %d\n", yellow, reset, s.Incons))
    sb.WriteString(fmt.Sprintf("%sFiles scanned:%s %d\n", green, reset, s.Files))
    sb.WriteString(fmt.Sprintf("%sTables skipped:%s %d\n", red, reset, s.Skipped))
    return sb.String()
}

func ProcessDir(root, outDir, delim string, skipCode bool) (*Summary, error) {
    var summary Summary
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        if strings.HasSuffix(d.Name(), ".md") {
            summary.Files++
            tables, err := processFile(path, outDir, delim, skipCode)
            if err != nil {
                return err
            }
            summary.Tables += tables.Tables
            summary.Rows += tables.Rows
            summary.Incons += tables.Incons
            summary.Skipped += tables.Skipped
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return &summary, nil
}

type TableStats struct {
    Tables  int
    Rows    int
    Incons  int
    Skipped int
}

func processFile(path, outDir, delim string, skipCode bool) (*TableStats, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    var lines []string
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    tables := parseTables(lines)
    stats := &TableStats{}
    for i, tbl := range tables {
        if skipCode && containsCodeOrImage(tbl) {
            stats.Skipped++
            continue
        }
        stats.Tables++
        stats.Rows += len(tbl)
        if inconsistent(tbl) {
            stats.Incons++
        }
        err := writeCSV(path, outDir, delim, i+1, tbl)
        if err != nil {
            return nil, err
        }
    }
    return stats, nil
}

func parseTables(lines []string) [][]string {
    var tables [][]string
    var current []string
    inTable := false
    for i := 0; i < len(lines); i++ {
        line := strings.TrimSpace(lines[i])
        if strings.Contains(line, "|") {
            if !inTable {
                // look ahead for separator
                if i+1 < len(lines) && strings.Contains(strings.TrimSpace(lines[i+1]), "|") {
                    sep := strings.TrimSpace(lines[i+1])
                    if strings.Contains(sep, "-") {
                        inTable = true
                        current = []string{line}
                        i++ // skip separator
                        continue
                    }
                }
            } else {
                current = append(current, line)
                continue
            }
        }
        if inTable {
            tables = append(tables, current)
            current = nil
            inTable = false
        }
    }
    if inTable && len(current) > 0 {
        tables = append(tables, current)
    }
    return tables
}

func containsCodeOrImage(tbl []string) bool {
    for _, row := range tbl {
        if strings.Contains(row, "`") || strings.Contains(row, "![") {
            return true
        }
    }
    return false
}

func inconsistent(tbl []string) bool {
    if len(tbl) == 0 {
        return false
    }
    cols := countPipes(tbl[0])
    for _, row := range tbl[1:] {
        if countPipes(row) != cols {
            return true
        }
    }
    return false
}

func countPipes(row string) int {
    parts := strings.Split(row, "|")
    if len(parts) > 0 && parts[0] == "" {
        parts = parts[1:]
    }
    if len(parts) > 0 && parts[len(parts)-1] == "" {
        parts = parts[:len(parts)-1]
    }
    return len(parts)
}

func writeCSV(srcPath, outDir, delim string, idx int, tbl []string) error {
    base := filepath.Base(srcPath)
    name := strings.TrimSuffix(base, filepath.Ext(base))
    outName := fmt.Sprintf("%s_table%d.csv", name, idx)
    if outDir != "" {
        outName = filepath.Join(outDir, outName)
    }
    f, err := os.Create(outName)
    if err != nil {
        return err
    }
    defer f.Close()
    writer := csv.NewWriter(f)
    writer.Comma = []rune(delim)[0]
    for _, row := range tbl {
        cells := splitRow(row)
        if err := writer.Write(cells); err != nil {
            return err
        }
    }
    writer.Flush()
    return writer.Error()
}

func splitRow(row string) []string {
    parts := strings.Split(row, "|")
    if len(parts) > 0 && parts[0] == "" {
        parts = parts[1:]
    }
    if len(parts) > 0 && parts[len(parts)-1] == "" {
        parts = parts[:len(parts)-1]
    }
    for i, p := range parts {
        parts[i] = strings.TrimSpace(p)
    }
    return parts
}
