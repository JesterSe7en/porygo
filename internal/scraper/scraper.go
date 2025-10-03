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

type Scraper struct {
	client *http.Client
	log    *logger.Logger
	cfg    *config.Config
	cache  storage.CacheStorage
}

// TODO: Look into goquery library to parse html better

func New(cfg *config.Config, log *logger.Logger, cache storage.CacheStorage) *Scraper {
	return &Scraper{
		client: new(http.Client),
		log:    log,
		cfg:    cfg,
		cache:  cache,
	}
}

// ScrapeWithRetry is the main public function that orchestrates scraping with caching and retry logic
func (s *Scraper) ScrapeWithRetry(url string) wp.Result {
	if !s.cfg.Force {
		if cached := s.checkCache(url); cached != nil {
			return *cached
		}
	}

	result := s.performScrapeWithRetries(url)

	if result.Err != nil {
		s.log.Error("Failed to scrape %s: %v", url, result.Err)
		return result
	}

	var data []byte
	switch v := result.Value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		s.log.Warn("Unsupported result type for %s: %T", url, v)
	}

	if len(data) > 0 {
		s.storeCacheResult(url, data)
	}

	return result
}

// performScrapeWithRetries handles the retry logic for scraping
func (s *Scraper) performScrapeWithRetries(url string) wp.Result {
	var lastErr error

	s.log.Debug("Starting scrape retry loop for URL %s with %d retries.", url, s.cfg.Retry)

	for attempt := 1; attempt <= s.cfg.Retry; attempt++ {
		s.log.Info("Attempting to scrape URL %s (attempt %d of %d)", url, attempt, s.cfg.Retry)

		result := s.scrape(url)
		if result.Err == nil {
			s.log.Info("Successfully scraped URL %s.", url)
			return result
		}

		lastErr = result.Err
		// Don't print out the stack trace
		s.log.Warn("Scraping attempt %d for URL %s failed: %s", attempt, url, result.Err.Error())

		// Wait before retry (except for last attempt)
		if attempt < s.cfg.Retry {
			delay := s.calculateBackoffDelay(attempt - 1)
			s.log.Info("Waiting %v before the next retry.", delay)
			time.Sleep(delay)
		}
	}

	return wp.Result{
		Value: nil,
		Err:   fmt.Errorf("all attempts failed: %s", lastErr.Error()),
	}
}

// scrape performs the actual HTTP request and returns the result
func (s *Scraper) scrape(url string) wp.Result {
	// cancel here allows us to gracefully clean up and stop the request if in the event
	// the user force quits e.g. ctrl-c
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
	defer cancel()

	start := time.Now()
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
	res, err := s.client.Do(req)
	if err != nil {
		return wp.Result{Value: nil, Err: err}
	}
	finished := time.Now()
	elapsed := finished.Sub(start)
	defer res.Body.Close()

	// Reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status#server_error_responses
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return wp.Result{Value: nil, Err: fmt.Errorf("request failed with status code: %d", res.StatusCode)}
	}

	data := ScrapedData{
		URL:          url,
		Status:       res.StatusCode,
		Title:        res.Header.Get("Title"),
		ContentType:  res.Header.Get("Content-Type"),
		Size:         res.ContentLength,
		ResponseTime: elapsed,
		Timestamp:    finished,
	}

	// and css selectors is given
	if data.ContentType == "text/html" && len(s.cfg.SelectorsConfig.Select) > 0 {
		// use the css selectors to scrape data from the HTML content
		extractBySelector(res, &data, s.cfg.SelectorsConfig.Select)
	}

	// if regex used, search the HTML content that was scrapped via CSS selectors
	if data.ContentType == "text/html" && len(s.cfg.SelectorsConfig.Pattern) > 0 {
		// use the regex patterns to scrape data from the HTML content
		filterByRegex(res, &data, s.cfg.SelectorsConfig.Pattern)
	}

	return wp.Result{
		Value: data,
		Err:   nil,
	}
}

func extractBySelector(res *http.Response, data *ScrapedData, regex []string) {
	if res == nil || len(regex) == 0 {
		return
	}
	panic("unimplemented")
}

func filterByRegex(res *http.Response, data *ScrapedData, selectors []string) {
	if res == nil || len(selectors) == 0 {
		return
	}
	panic("unimplemented")
}

// calculateBackoffDelay calculates the delay for exponential backoff with optional jitter
// Uses fixed multiplier of 2.0 and auto-calculated max delay for production safety
func (s *Scraper) calculateBackoffDelay(attempt int) time.Duration {
	const multiplier = 2.0

	// Calculate exponential backoff: baseDelay * 2^attempt
	// Formula from: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
	delay := float64(s.cfg.Backoff.BaseDelay) * math.Pow(multiplier, float64(attempt))

	// Auto-calculate reasonable max delay: max(30s, baseDelay * 16)
	maxDelay := 30 * time.Second
	if baseMax := s.cfg.Backoff.BaseDelay * 16; baseMax > maxDelay {
		maxDelay = baseMax
	}

	if time.Duration(delay) > maxDelay {
		delay = float64(maxDelay)
	}

	if s.cfg.Backoff.Jitter {
		delay = rand.Float64() * delay
	}

	return time.Duration(delay)
}

// checkCache retrieves and validates cached data for the given URL
func (s *Scraper) checkCache(url string) *wp.Result {
	cached, err := s.cache.Get(context.Background(), url)

	if err != nil && err != storage.ErrNotFound {
		s.log.Error("Failed to retrieve %s from cache: %v", url, err)
		return nil
	}

	if err == storage.ErrNotFound {
		return nil
	}

	if time.Now().After(cached.ExpirationTime) {
		s.cleanupExpiredCache(url)
		return nil
	}

	// Return valid cached result
	s.log.Debug("Using cached data for %s, not expired yet", url)
	return &wp.Result{
		Value: cached,
		Err:   nil,
	}
}

// storeCacheResult stores the scraped result in the cache
func (s *Scraper) storeCacheResult(url string, data []byte) {
	s.log.Debug("Adding %s to cache...", url)

	entry := storage.CacheEntry{
		ExpirationTime: time.Now().Add(s.cfg.Database.Expiration),
		Value:          data,
	}

	err := s.cache.Set(context.Background(), url, entry)
	if err != nil {
		s.log.Error("Failed to store %s in cache: %v", url, err)
	}

	s.log.Debug("Cache put operation successful.")
}

// cleanupExpiredCache removes expired cache entries
func (s *Scraper) cleanupExpiredCache(url string) {
	s.log.Debug("Cached data for %s is old, discarding...", url)
	if err := s.cache.Delete(context.Background(), url); err != nil {
		s.log.Error("Failed to delete %s from cache: %v", url, err)
	}
}
