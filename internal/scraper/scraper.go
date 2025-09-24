// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package scraper involves all functions related to actually doing http requests and scraping the
// data from the response
package scraper

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JesterSe7en/scrapego/internal/database"
	"github.com/JesterSe7en/scrapego/internal/logger"
	wp "github.com/JesterSe7en/scrapego/internal/workerpool"
)

// TODO: Look into goquery library to parse html better

var client = &http.Client{}

func ScrapeWithRetry(url string, timeout time.Duration, retries int, backoff time.Duration) wp.Result {

	// check cache if available
	cachedEntry, err := database.Get([]byte(url))
	found := false
	// check if the error is something else other than not found.
	// if err is not found, don't execute the cache check
	if err != nil && err != database.ErrNotFound {
		logger.Error("Failed to retrieve %s from cache: %v", url, err)
	} else {
		if err == nil {
			found = true
		}
	}

	if found {
		var expiration = time.Hour
		// check if data is more than an hour old
		entryTime := time.Unix(cachedEntry.Timestamp, 0)
		if time.Since(entryTime) < expiration {
			logger.Debug("Using cached data for %s, not expired yet", url)
			return wp.Result{
				Value: cachedEntry.Data,
				Err:   nil,
			}
		} else {
			logger.Debug("Cached data for %s is old, discarding...", url)
			err := database.Delete([]byte(url))
			if err != nil {
				logger.Error("Failed to delete %s from cache: %v", url, err)
			}
		}
	}

	var lastErr error
	logger.Debug("Starting scrape retry loop for URL %s with %d retries.", url, retries)
	for attempt := 1; attempt <= retries; attempt++ {
		logger.Info("Attempting to scrape URL %s (attempt %d of %d)", url, attempt, retries)
		result := scrape(url, timeout)
		if result.Err == nil {

			logger.Info("Successfully scraped URL %s.", url)
			// add to cache
			logger.Debug("Adding %s to cache...", url)
			database.Put([]byte(url), []byte("this is the cache for "+url))

			logger.Debug("Cache put operation successful.")
			return result
		}
		lastErr = result.Err

		logger.Warn("Scraping attempt %d for URL %s failed: %v", attempt, url, result.Err)
		if attempt < retries {
			logger.Info("Waiting %v before the next retry.", backoff)
			time.Sleep(backoff)
		}
	}

	return wp.Result{
		Value: nil,
		Err:   fmt.Errorf("all attempts failed: %s", lastErr.Error()),
	}
}

func scrape(url string, timeout time.Duration) wp.Result {
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

	// Reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status#server_error_responses
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
