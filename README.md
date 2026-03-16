# go-batch-dl

A simple, concurrent batch downloader written in Go. It scrapes a target URL for links matching a specific extension and downloads them in parallel.

## Build and Run

You can run the tool directly using `go run`:

```bash
go run cmd/gobatchdl/main.go -url "http://example.com" -ext ".jpg"
```

## Options

- `-url`: (Required) The target URL to scrape.
- `-ext`: The file extension filter (default: ".jpg").
- `-dir`: The destination directory (default: "./downloads").
- `-workers`: The number of concurrent download workers (default: 5).

## Structure

- `cmd/gobatchdl`: Main entry point.
- `internal/downloader`: HTTP fetching and worker pool implementation.
- `internal/scraper`: Link extraction logic.
