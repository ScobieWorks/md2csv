package main

import (
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestParseTables(t *testing.T) {
    content := `
# Sample

| Header1 | Header2 |
|---------|---------|
| Cell1   | Cell2   |
| Cell3   | Cell4   |

Some text

| A | B | C |
|---|---|---|
| 1 | 2 | 3 |
`
    lines := strings.Split(content, "\n")
    tables := parseTables(lines)
    if len(tables) != 2 {
        t.Fatalf("expected 2 tables, got %d", len(tables))
    }
    if len(tables[0]) != 3 {
        t.Fatalf("first table rows expected 3, got %d", len(tables[0]))
    }
}

func TestWriteCSV(t *testing.T) {
    tmpDir := t.TempDir()
    src := filepath.Join(tmpDir, "test.md")
    content := `
| H1 | H2 |
|----|----|
| a  | b  |
| c  | d  |
`
    if err := ioutil.WriteFile(src, []byte(content), 0644); err != nil {
        t.Fatalf("write file: %v", err)
    }
    tbl := parseTables(strings.Split(content, "\n"))[0]
    if err := writeCSV(src, tmpDir, ",", 1, tbl); err != nil {
        t.Fatalf("writeCSV: %v", err)
    }
    csvPath := filepath.Join(tmpDir, "test_table1.csv")
    data, err := os.ReadFile(csvPath)
    if err != nil {
        t.Fatalf("read csv: %v", err)
    }
    expected := "H1,H2\na,b\nc,d\n"
    if string(data) != expected {
        t.Fatalf("csv content mismatch\nexpected:\n%s\ngot:\n%s", expected, string(data))
    }
}
