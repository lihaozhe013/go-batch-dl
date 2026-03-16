package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lihaozhe013/go-batch-dl/internal/downloader"
	"github.com/lihaozhe013/go-batch-dl/internal/scraper"
)

func main() {
	// ==========================================
	// 1. CLI Arguments Parsing
	// ==========================================
	var (
		targetURL string
		extFilter string
		destDir   string
		workerNum int
	)

	flag.StringVar(&targetURL, "url", "", "Target URL")
	flag.StringVar(&extFilter, "ext", ".jpg", "File extension filter (e.g. .jpg, .png)")
	flag.StringVar(&destDir, "dir", "./downloads", "Destination directory")
	flag.IntVar(&workerNum, "workers", 5, "Number of concurrent workers")
	flag.Parse()

	if targetURL == "" {
		fmt.Println("Please verify the target URL using -url")
		return
	}

	// Print configuration for debugging
	fmt.Println("Scanning...")
	fmt.Printf("Target: %s, Ext: %s, Dir: %s, Workers: %d\n", targetURL, extFilter, destDir, workerNum)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	// ==========================================
	// 2. Single Thread: Fetch HTML and Extract Links
	// ==========================================
	htmlData, err := downloader.FetchHTML(targetURL)
	if err != nil {
		fmt.Printf("Failed to fetch HTML: %v\n", err)
		return
	}

	links, err := scraper.ExtractLinks(htmlData, targetURL, extFilter)
	if err != nil {
		fmt.Printf("Failed to extract links: %v\n", err)
		return
	}

	if len(links) == 0 {
		fmt.Println("No matching files found")
		return
	}

	fmt.Printf("Found %d links, starting download...\n", len(links))

	// ==========================================
	// 3. Multi-threading: Initialize Worker Pool Components
	// ==========================================
	jobs := make(chan downloader.DownloadJob, len(links))
	results := make(chan downloader.DownloadResult, len(links))
	var wg sync.WaitGroup

	// ==========================================
	// 4. Start Workers (Concurrency Core)
	// ==========================================
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go downloader.Worker(i, &wg, jobs, results)
	}

	// ==========================================
	// 5. Dispatch Jobs and Collect Results
	// ==========================================
	go func() {
		for _, link := range links {
			filename := filepath.Base(link)
			// Remove query parameters if present
			if idx := strings.Index(filename, "?"); idx != -1 {
				filename = filename[:idx]
			}
			// Just in case filename is empty or invalid
			if filename == "" || filename == "." || filename == "/" {
				filename = fmt.Sprintf("file_%d%s", len(link), extFilter)
			}

			jobs <- downloader.DownloadJob{
				URL:      link,
				DestPath: filepath.Join(destDir, filename),
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	successCount := 0
	failCount := 0
	for result := range results {
		if result.Error != nil {
			fmt.Printf("[FAIL] %s: %v\n", result.URL, result.Error)
			failCount++
		} else {
			fmt.Printf("[OK] %s\n", result.URL)
			successCount++
		}
	}

	fmt.Printf("\nDownload completed! Success: %d, Fail: %d\n", successCount, failCount)
}
