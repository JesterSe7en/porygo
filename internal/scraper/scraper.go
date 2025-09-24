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

// ScrapeWithRetry is the main public function that orchestrates scraping with caching and retry logic
func ScrapeWithRetry(url string, timeout time.Duration, retries int, backoff time.Duration) wp.Result {
	if cachedResult := checkCache(url); cachedResult != nil {
		return *cachedResult
	}

	result := performScrapeWithRetries(url, timeout, retries, backoff)

	if result.Err != nil {
		logger.Error("failed to scrape %s: %v", url, result.Err)
	} else {
		storeCacheResult(url)
	}
	return result
}

// scrape performs the actual HTTP request and returns the result
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

// performScrapeWithRetries handles the retry logic for scraping
func performScrapeWithRetries(url string, timeout time.Duration, retries int, backoff time.Duration) wp.Result {
	var lastErr error

	logger.Debug("Starting scrape retry loop for URL %s with %d retries.", url, retries)

	for attempt := 1; attempt <= retries; attempt++ {
		logger.Info("Attempting to scrape URL %s (attempt %d of %d)", url, attempt, retries)

		result := scrape(url, timeout)
		if result.Err == nil {
			logger.Info("Successfully scraped URL %s.", url)
			return result
		}

		lastErr = result.Err
		logger.Warn("Scraping attempt %d for URL %s failed: %v", attempt, url, result.Err)

		// Wait before retry (except for last attempt)
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

// checkCache retrieves and validates cached data for the given URL
func checkCache(url string) *wp.Result {
	cachedEntry, err := database.Get([]byte(url))

	// Handle database errors (not including "not found")
	if err != nil && err != database.ErrNotFound {
		logger.Error("Failed to retrieve %s from cache: %v", url, err)
		return nil
	}

	// No cached entry found
	if err == database.ErrNotFound {
		return nil
	}

	// Check if cached entry is expired
	if isCacheExpired(&cachedEntry) {
		cleanupExpiredCache(url)
		return nil
	}

	// Return valid cached result
	logger.Debug("Using cached data for %s, not expired yet", url)
	return &wp.Result{
		Value: cachedEntry.Data,
		Err:   nil,
	}
}

// isCacheExpired checks if a cache entry has exceeded its expiration time
// TODO: not sure if I should pass pointer here. Data could be large theoretically
func isCacheExpired(entry *database.CacheEntry) bool {
	const expiration = time.Hour
	entryTime := time.Unix(entry.Timestamp, 0)
	return time.Since(entryTime) >= expiration
}

// cleanupExpiredCache removes expired cache entries
func cleanupExpiredCache(url string) {
	logger.Debug("Cached data for %s is old, discarding...", url)
	if err := database.Delete([]byte(url)); err != nil {
		logger.Error("Failed to delete %s from cache: %v", url, err)
	}
}

// storeCacheResult stores the scraped result in the cache
func storeCacheResult(url string) {
	logger.Debug("Adding %s to cache...", url)
	// TODO: actually put real data into the cache
	database.Put([]byte(url), []byte("this is the cache for "+url))
	logger.Debug("Cache put operation successful.")
}
