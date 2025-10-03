package scraper

import (
	"time"
)

type ScrapedData struct {
	URL          string        `json:"url"`
	Status       int           `json:"status"`
	Title        string        `json:"title,omitempty"`
	ContentType  string        `json:"content_type,omitempty"`
	Size         int64         `json:"size,omitempty"`
	ResponseTime time.Duration `json:"response_time"`
	Timestamp    time.Time     `json:"timestamp"`

	// CSS selector results
	Extracted map[string][]string `json:"extracted,omitempty"`

	// Regex matches
	Matches map[string][]string `json:"matches,omitempty"`
}
