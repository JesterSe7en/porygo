// Package scraper involves all functions related to actually doing http requests and scraping the
// data from the response
package scraper

import (
	"context"
	"fmt"
	"net/http"
	"time"

	wp "github.com/JesterSe7en/scrapego/internal/workerpool"
)

// TODO: Look into goquery library to parse html better

var client = &http.Client{}

func ScrapeWithTimeout(url string, timeout time.Duration) wp.Result {
	// cancel here allows us to gracefully clean up and stop the request if in the event
	// the user force quits e.g. ctrl-c
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return wp.Result{
			Value: nil,
			Err:   err,
		}
	}

	// Set the User-Agent header to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	// Execute the request
	res, err := client.Do(req)
	if err != nil {
		return wp.Result{Value: nil, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return wp.Result{Value: nil, Err: fmt.Errorf("request failed with status code: %d", res.StatusCode)}
	}

	// Extract content type
	contentType := res.Header.Get("Content-Type")

	return wp.Result{
		Value: contentType,
		Err:   nil,
	}
}
