package downloader

import (
	"fmt"
	"io"
	"net/http"
)

// FetchHTML initiates an HTTP GET request and returns the page content as a string
// Note: The function name must be capitalized to be called by other packages (e.g. main)
func FetchHTML(targetURL string) (string, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
