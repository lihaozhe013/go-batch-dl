package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type DownloadJob struct {
	URL      string
	DestPath string
}

type DownloadResult struct {
	URL   string
	Error error
}

// a Worker is the code that each concurrent coroutine executes
func Worker(id int, wg *sync.WaitGroup, jobs <-chan DownloadJob, results chan<- DownloadResult) {
	defer wg.Done()
	for job := range jobs {
		err := downloadFile(job.URL, job.DestPath)
		results <- DownloadResult{
			URL:   job.URL,
			Error: err,
		}
	}
}

func downloadFile(url string, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
