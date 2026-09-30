# md2csv

**Markdown Table to CSV Converter**

`md2csv` is a small, offline command‑line tool written in Go that scans a directory for Markdown files, extracts all tables, and writes each table to a separate CSV file. It preserves cell content, trims whitespace, handles multiline cells, and reports any rows with mismatched column counts.

## Features

- Recursively search for `*.md` files.
- Export each table as a CSV file named `<source>_table<index>.csv`.
- Custom delimiter: comma (default), tab, or pipe.
- Optional flag to skip tables containing code blocks or images.
- Color‑coded summary report (tables processed, rows, inconsistencies, files scanned, tables skipped).

## Installation

```bash
# From source
go install github.com/ScobieWorks/md2csv@latest
```

## Usage

```bash
# Basic usage
md2csv /path/to/docs

# Specify output directory
md2csv -o ./csvs /path/to/docs

# Use tab delimiter
md2csv -d "\t" /path/to/docs

# Skip tables with code blocks or images
md2csv -s /path/to/docs
```

## Testing

```bash
go test ./...
```

## License

MIT


## Support

If this project saved you time, optional support is welcome: https://paypal.me/Damonwill
