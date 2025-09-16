// Package scraper involves all functions related to actually doing http requests and scraping the
// data from the response
package scraper

import (
	"context"
	"errors"
	"net/http"
	"time"

	wp "github.com/JesterSe7en/scrapego/internal/workerpool"
)

func ScrapeWithTimeout(url string, timeout time.Duration) wp.Result {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return wp.Result{
			Value: nil,
			Err:   err,
		}
	}

	// Create a client
	client := &http.Client{}

	// Execute the request
	res, err := client.Do(req)
	if err != nil {
		return wp.Result{Value: nil, Err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return wp.Result{Value: nil, Err: errors.New("did not recieve status code 200")}
	}

	// Extract content type
	contentType := res.Header.Get("Content-Type")

	return wp.Result{
		Value: contentType,
		Err:   nil,
	}
}
