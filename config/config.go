// Package config contains the key-value pairs that deal with how to use the scrapego tool
package config

import "time"

type Config struct {
	Input       string
	Concurrency int
	Timeout     time.Duration
	Output      string
	Verbose     bool
	Retry       int
	Rate        int
	Force       bool
}

//  future add maybe a way to marshall this into a yaml or something else
