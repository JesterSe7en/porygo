// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package scraper involves all functions related to actually doing http requests and scraping the
// data from the response
package scraper

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/JesterSe7en/scrapego/config"
	"github.com/JesterSe7en/scrapego/internal/logger"
	"github.com/JesterSe7en/scrapego/internal/storage"
	wp "github.com/JesterSe7en/scrapego/internal/workerpool"
)

// TODO: Look into goquery library to parse html better
var client = &http.Client{}

// ScrapeWithRetry is the main public function that orchestrates scraping with caching and retry logic
func ScrapeWithRetry(url string, retries int, timeout time.Duration, backoffConfig config.BackoffConfig, expiration time.Duration, cache storage.CacheStorage) wp.Result {
	if cachedResult := checkCache(url, cache); cachedResult != nil {
		return *cachedResult
	}

	result := performScrapeWithRetries(url, timeout, retries, backoffConfig)

	if result.Err != nil {
		logger.Error("Failed to scrape %s: %v", url, result.Err)
	} else {
		var data []byte
		switch v := result.Value.(type) {
		case []byte:
			data = v
		case string:
			data = []byte(v)
		}

		storeCacheResult(url, data, expiration, cache)
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
func performScrapeWithRetries(url string, timeout time.Duration, retries int, backoffConfig config.BackoffConfig) wp.Result {
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
		// Don't print out the stack trace
		logger.Warn("Scraping attempt %d for URL %s failed: %s", attempt, url, result.Err.Error())

		// Wait before retry (except for last attempt)
		if attempt < retries {
			delay := calculateBackoffDelay(attempt-1, backoffConfig)
			logger.Info("Waiting %v before the next retry.", delay)
			time.Sleep(delay)
		}
	}

	return wp.Result{
		Value: nil,
		Err:   fmt.Errorf("all attempts failed: %s", lastErr.Error()),
	}
}

// checkCache retrieves and validates cached data for the given URL
func checkCache(url string, cache storage.CacheStorage) *wp.Result {
	cached, err := cache.Get(context.Background(), url)

	if err != nil && err != storage.ErrNotFound {
		logger.Error("Failed to retrieve %s from cache: %v", url, err)
		return nil
	}

	if err == storage.ErrNotFound {
		return nil
	}

	if time.Now().After(cached.ExpirationTime) {
		cleanupExpiredCache(url, cache)
		return nil
	}

	// Return valid cached result
	logger.Debug("Using cached data for %s, not expired yet", url)
	return &wp.Result{
		Value: cached,
		Err:   nil,
	}
}

// cleanupExpiredCache removes expired cache entries
func cleanupExpiredCache(url string, cache storage.CacheStorage) {
	logger.Debug("Cached data for %s is old, discarding...", url)
	if err := cache.Delete(context.Background(), url); err != nil {
		logger.Error("Failed to delete %s from cache: %v", url, err)
	}
}

// storeCacheResult stores the scraped result in the cache
func storeCacheResult(url string, data []byte, expiration time.Duration, cache storage.CacheStorage) {
	logger.Debug("Adding %s to cache...", url)

	entry := storage.CacheEntry{
		ExpirationTime: time.Now().Add(expiration),
		Value:          data,
	}

	err := cache.Set(context.Background(), url, entry)
	if err != nil {
		logger.Error("Failed to store %s in cache: %v", url, err)
	}

	logger.Debug("Cache put operation successful.")
}

// calculateBackoffDelay calculates the delay for exponential backoff with optional jitter
// Uses fixed multiplier of 2.0 and auto-calculated max delay for production safety
func calculateBackoffDelay(attempt int, cfg config.BackoffConfig) time.Duration {
	const multiplier = 2.0

	// Calculate exponential backoff: baseDelay * 2^attempt
	// Formula from: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
	delay := float64(cfg.BaseDelay) * math.Pow(multiplier, float64(attempt))

	// Auto-calculate reasonable max delay: max(30s, baseDelay * 16)
	maxDelay := 30 * time.Second
	if baseMax := cfg.BaseDelay * 16; baseMax > maxDelay {
		maxDelay = baseMax
	}

	if time.Duration(delay) > maxDelay {
		delay = float64(maxDelay)
	}

	if cfg.Jitter {
		delay = rand.Float64() * delay
	}

	return time.Duration(delay)
}
