package util

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
)

func extractWgArticleID(scriptContent string) string {
	// Regular expression to match "wgArticleId" and its value
	re := regexp.MustCompile(`"wgArticleId":(\d+)`)

	// Find matches
	matches := re.FindStringSubmatch(scriptContent)

	// If there is a match, return the value
	if len(matches) == 2 {
		return matches[1]
	}

	// Otherwise, return an empty string
	return ""
}

func URLToPageID(url string) (string, error) {
	c := colly.NewCollector()

	var wgArticleID string
	var err error
	stopScraping := false // Flag to indicate whether scraping should stop

	// Visit the page
	c.OnHTML("script", func(e *colly.HTMLElement) {
		if stopScraping {
			return // Stop scraping if the flag is set
		}

		// Extract the content of the script tag
		scriptContent := e.Text

		// Find and extract the value of "wgArticleId"
		wgArticleID = extractWgArticleID(scriptContent)
		if wgArticleID != "" {
			stopScraping = true // Set the flag to stop scraping
		}
	})

	// Error handling for scraping
	c.OnError(func(r *colly.Response, e error) {
		err = fmt.Errorf("failed to scrape URL: %s, error: %v", r.Request.URL, e)
	})

	c.Visit(url)

	// Return the values after scraping is complete
	return wgArticleID, err
}
