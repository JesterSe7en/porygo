// Package config contains the key-value pairs that deal with how to use the scrapego tool
package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/JesterSe7en/scrapego/internal/logger"
)

// Config holds all configuration options for the scrapego tool
type Config struct {
	Input       string        `toml:"input"`
	Concurrency int           `toml:"concurrency"`
	Timeout     time.Duration `toml:"timeout"`
	Output      string        `toml:"output"`
	Verbose     bool          `toml:"verbose"`
	Retry       int           `toml:"retry"`
	Backoff     time.Duration `toml:"backoff"`
	Force       bool          `toml:"force"`
}

// Manager handles configuration loading, merging, and saving
type Manager struct {
	configPath string
}

// NewManager creates a new configuration manager
func NewManager(configPath string) *Manager {
	return &Manager{
		configPath: configPath,
	}
}

// DefaultManager creates a configuration manager with the default config path
func DefaultManager() *Manager {
	return NewManager("config.toml")
}

// Defaults returns a Config struct with default values
func Defaults() Config {
	return Config{
		Input:       "",
		Concurrency: 5,
		Timeout:     10 * time.Second,
		Output:      "JSON",
		Verbose:     false,
		Retry:       3,
		Backoff:     2 * time.Second,
		Force:       false,
	}
}

// Load loads configuration with the following precedence:
// 1. Default values
// 2. Config file (if exists)
// 3. Environment variables (TODO)
// 4. CLI flags (handled by caller)
func (m *Manager) Load() (Config, error) {
	// Start with defaults
	cfg := Defaults()

	// Try to load from file if it exists
	logger.Debug("Attempting to load config file, uses default if none found...")
	if _, err := os.Stat(m.configPath); err == nil {
		logger.Debug("Found config file")
		fileCfg, err := m.loadFromFile()
		if err != nil {
			return cfg, fmt.Errorf("failed to load config file: %s", err.Error())
		}
		cfg = m.mergeConfigs(cfg, fileCfg)
	}

	return cfg, nil
}

// loadFromFile loads configuration from a TOML file
func (m *Manager) loadFromFile() (Config, error) {
	var cfg Config

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %s", err.Error())
	}

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %s", err.Error())
	}

	return cfg, nil
}

// mergeConfigs merges two config structs, with the second taking precedence
func (m *Manager) mergeConfigs(base, override Config) Config {
	result := base

	// Only override non-zero values
	if override.Input != "" {
		result.Input = override.Input
	}
	if override.Concurrency != 0 {
		result.Concurrency = override.Concurrency
	}
	if override.Timeout != 0 {
		result.Timeout = override.Timeout
	}
	if override.Output != "" {
		result.Output = override.Output
	}
	// Booleans are trickier - we assume false is intentional in config files
	result.Verbose = override.Verbose
	result.Force = override.Force

	if override.Retry != 0 {
		result.Retry = override.Retry
	}
	if override.Backoff != 0 {
		result.Backoff = override.Backoff
	}

	return result
}

// Save writes the configuration to a TOML file
func (m *Manager) Save(cfg Config) error {
	return m.SaveWithForce(cfg, false)
}

// SaveWithForce writes the configuration to a TOML file
// If force is true, it will overwrite an existing config file
func (m *Manager) SaveWithForce(cfg Config, force bool) error {
	// Check if config file already exists
	if !force {
		if _, err := os.Stat(m.configPath); err == nil {
			return fmt.Errorf("config file %s already exists. Use force to overwrite", m.configPath)
		}
	}

	// Ensure the directory exists
	if dir := filepath.Dir(m.configPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create config directory: %s", err.Error())
		}
	}

	buffer, err := m.encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to encode config: %s", err.Error())
	}

	err = os.WriteFile(m.configPath, buffer.Bytes(), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %s", err.Error())
	}

	return nil
}

// InitDefaults creates a config file with default values
func (m *Manager) InitDefaults() error {
	return m.InitDefaultsWithForce(false)
}

// InitDefaultsWithForce creates a config file with default values
// If force is true, it will overwrite an existing config file
func (m *Manager) InitDefaultsWithForce(force bool) error {
	defaults := Defaults()
	return m.SaveWithForce(defaults, force)
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

	if cfg.Output != "JSON" && cfg.Output != "CSV" {
		errs = append(errs, "output must be either 'JSON' or 'CSV'")
	}

	if cfg.Retry < 0 {
		errs = append(errs, "retry count cannot be negative")
	}

	if cfg.Backoff <= 0 {
		errs = append(errs, "backoff must be greater than 0")
	}

	if len(errs) > 0 {
		return errors.New("configuration validation failed: " + fmt.Sprintf("%v", errs))
	}

	return nil
}

// String returns a string representation of the config
func (cfg *Config) String() string {
	var buffer bytes.Buffer
	toml.NewEncoder(&buffer).Encode(cfg)
	return buffer.String()
}

// Deprecated functions for backward compatibility
// These will be removed in a future version

// DefaultConfig returns default configuration values
// Deprecated: Use Defaults() instead
func DefaultConfig() Config {
	return Defaults()
}

// Encode converts config to TOML buffer
// Deprecated: Use Manager.encode() instead
func Encode(cfg *Config) (bytes.Buffer, error) {
	if cfg == nil {
		return bytes.Buffer{}, errors.New("config is nil")
	}

	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(cfg)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("cannot encode config: %s", err.Error())
	}

	return buffer, nil
}

// WriteToToml writes config buffer to file
// Deprecated: Use Manager.Save() instead
func WriteToToml(buffer bytes.Buffer) error {
	return os.WriteFile("config.toml", buffer.Bytes(), 0o644)
}

// InitConfigFile creates config with defaults
// Deprecated: Use Manager.InitDefaults() instead
func InitConfigFile() error {
	return InitConfigFileWithForce(false)
}

// InitConfigFileWithForce creates config with defaults and force option
// Deprecated: Use Manager.InitDefaultsWithForce() instead
func InitConfigFileWithForce(force bool) error {
	manager := DefaultManager()
	return manager.InitDefaultsWithForce(force)
}
