package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ExtractLinks(htmlData string, baseURL string, extFilter string) ([]string, error) {
	var validURLs []string
	re := regexp.MustCompile(`href=["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(htmlData, -1)

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		link := match[1]

		// ignore some common non-file links
		if link == "" || link == "/" || link == "#" || strings.HasPrefix(link, "javascript:") || strings.HasPrefix(link, "mailto:") {
			continue
		}

		// check whether the extension matches extFilter (if extFilter is not empty)
		if extFilter != "" && !strings.HasSuffix(strings.ToLower(link), strings.ToLower(extFilter)) {
			continue
		}

		// if it is a relative path, concatenate it with the baseURL to form the complete download link
		u, err := url.Parse(link)
		if err != nil {
			continue
		}

		absoluteURL := base.ResolveReference(u).String()

		// remove duplicates
		if !seen[absoluteURL] {
			validURLs = append(validURLs, absoluteURL)
			seen[absoluteURL] = true
		}
	}

	return validURLs, nil
}
