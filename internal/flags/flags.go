// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package flags defines string constants for command-line flag names used throughout
// the application.
package flags

const (
	FlagLog         = "log"          // path to log file
	FlagDebug       = "debug"        // enable debug mode
	FlagConfig      = "config"       // path to config file
	FlagConcurrency = "concurrency"  // number of concurrent requests
	FlagTimeout     = "timeout"      // timeout duration for requests
	FlagVerbose     = "verbose"      // enable verbose mode
	FlagRetry       = "retry"        // number of retries for failed requests
	FlagRetryDelay  = "retry-delay"  // delay duration between retries
	FlagRetryJitter = "retry-jitter" // enable jitter for retry delays
	FlagBackoff     = "backoff"      // backoff duration between retries
	FlagForce       = "force"        // ignore cache and scrape fresh data

	// Scraper flags
	FlagSelect  = "select"  // CSS selectors
	FlagPattern = "pattern" // regex FlagPattern
	FlagFormat  = "format"  // output format json|csv|plain
	FlagQuiet   = "quiet"   // only output extracted data
	FlagHeaders = "headers" // include response headers
)
