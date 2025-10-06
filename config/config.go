// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package config contains the key-value pairs that deal with how to use the porygo tool
package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	configWriteMode = 0o644 // root r/w else r
	configDirMode   = 0o755 // root r/w/x else r/x
)

type Database struct {
	Expiration time.Duration `toml:"expiration"`
}

// BackoffConfig defines exponential backoff configuration
type BackoffConfig struct {
	BaseDelay time.Duration `toml:"base_delay"` // Initial delay between retries
	Jitter    bool          `toml:"jitter"`     // Whether to add jitter (default: true)
}

type SelectorsConfig struct {
	Select  []string `toml:"select"`  // css selectors
	Pattern []string `toml:"pattern"` // regex patterns
}

// Config holds all configuration options for the porygo tool
type Config struct {
	Concurrency     int             `toml:"concurrency"` // number of concurrent requests
	Timeout         time.Duration   `toml:"timeout"`     // timeout for each request
	Format          string          `toml:"format"`      // output format for the scraped data
	Retry           int             `toml:"retry"`       // number of retries for failed requests
	Backoff         BackoffConfig   `toml:"backoff"`     // exponential backoff configuration
	SelectorsConfig SelectorsConfig `toml:"selectors"`   // css/regex selectors configuration
	Database        Database        `toml:"database"`    // database configuration
	Force           bool            `toml:"force"`       // force scraping even if data exists
	Quiet           bool            `toml:"quiet"`       // suppress output, only show scrapped data
	Headers         bool            `toml:"headers"`     // include headers in output
}

type Manager struct {
	configPath string
}

func NewManager(configPath string) *Manager {
	return &Manager{
		configPath: configPath,
	}
}

func DefaultManager() *Manager {
	return NewManager("config.toml")
}

func Defaults() Config {
	return Config{
		Concurrency: 5,
		Timeout:     10 * time.Second,
		Format:      "json",
		Retry:       3,
		Backoff: BackoffConfig{
			BaseDelay: 1 * time.Second,
			Jitter:    true,
		},
		Quiet:   false,
		Headers: false,
		SelectorsConfig: SelectorsConfig{
			Select:  []string{},
			Pattern: []string{},
		},
		Force: false,
		Database: Database{
			Expiration: 24 * time.Hour,
		},
	}
}

func (m *Manager) LoadDefaults() Config {
	// Start with defaults
	return Defaults()
}

func (m *Manager) LoadFromFile(filePath string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %s", err.Error())
	}

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %s", err.Error())
	}

	return cfg, nil
}

func (m *Manager) Save(cfg Config) error {
	// Check if config file already exists
	// Ensure the directory exists
	if dir := filepath.Dir(m.configPath); dir != "." {
		if err := os.MkdirAll(dir, configDirMode); err != nil {
			return fmt.Errorf("failed to create config directory: %s", err.Error())
		}
	}

	buffer, err := m.encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to encode config: %s", err.Error())
	}

	err = os.WriteFile(m.configPath, buffer.Bytes(), configWriteMode)
	if err != nil {
		return fmt.Errorf("failed to write config file: %s", err.Error())
	}

	return nil
}

// InitDefaultsWithForce creates a config file with default values
// If force is true, it will overwrite an existing config file
func (m *Manager) InitDefaults() error {
	defaults := Defaults()
	return m.Save(defaults)
}

// encode converts the config into a TOML buffer
func (m *Manager) encode(cfg Config) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	err := toml.NewEncoder(&buffer).Encode(cfg)
	if err != nil {
		return buffer, fmt.Errorf("failed to encode config to TOML: %s", err.Error())
	}

	return buffer, nil
}

// Validate checks if the configuration values are valid
func (cfg *Config) Validate() error {
	var errs []string

	if cfg.Concurrency <= 0 {
		errs = append(errs, "concurrency must be greater than 0")
	}

	if cfg.Timeout <= 0 {
		errs = append(errs, "timeout must be greater than 0")
	}

	format := strings.ToLower(cfg.Format)
	if format != "json" && format != "text" {
		errs = append(errs, "format must be either 'json' or 'text'")
	}

	if cfg.Retry < 0 {
		errs = append(errs, "retry count cannot be negative")
	}

	if cfg.Backoff.BaseDelay <= 0 {
		errs = append(errs, "backoff base_delay must be greater than 0")
	}

	if len(errs) > 0 {
		return errors.New("configuration validation failed: " + strings.Join(errs, ", "))
	}

	return nil
}

// String returns a string representation of the config
func (cfg *Config) String() string {
	var buffer bytes.Buffer
	toml.NewEncoder(&buffer).Encode(cfg)
	return buffer.String()
}
