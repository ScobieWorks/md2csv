package main

import (
    "flag"
    "fmt"
    "log"
)

var (
    outputDir string
    delimiter string
    skipCode  bool
)

func init() {
    flag.StringVar(&outputDir, "o", "", "Output directory for CSV files (default: same as source)")
    flag.StringVar(&delimiter, "d", ",", "Delimiter for CSV output (comma, tab, pipe)")
    flag.BoolVar(&skipCode, "s", false, "Skip tables containing code blocks or images")
}

func main() {
    flag.Parse()
    if flag.NArg() == 0 {
        log.Fatal("Usage: md2csv <path>")
    }
    root := flag.Arg(0)
    summary, err := ProcessDir(root, outputDir, delimiter, skipCode)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Println(summary.Report())
}
