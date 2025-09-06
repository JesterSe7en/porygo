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
)

// Config holds all configuration options for the scrapego tool
type Config struct {
	Input       string        `toml:"input"`
	Concurrency int           `toml:"concurrency"`
	Timeout     time.Duration `toml:"timeout"`
	Output      string        `toml:"output"`
	Verbose     bool          `toml:"verbose"`
	Retry       int           `toml:"retry"`
	Rate        int           `toml:"rate"`
	Force       bool          `toml:"force"`
}

// ConfigManager handles configuration loading, merging, and saving
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		configPath: configPath,
	}
}

// DefaultConfigManager creates a configuration manager with the default config path
func DefaultConfigManager() *ConfigManager {
	return NewConfigManager("config.toml")
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
		Rate:        1,
		Force:       false,
	}
}

// Load loads configuration with the following precedence:
// 1. Default values
// 2. Config file (if exists)
// 3. Environment variables (TODO)
// 4. CLI flags (handled by caller)
func (m *ConfigManager) Load() (Config, error) {
	// Start with defaults
	cfg := Defaults()

	// Try to load from file if it exists
	if _, err := os.Stat(m.configPath); err == nil {
		fileCfg, err := m.loadFromFile()
		if err != nil {
			return cfg, fmt.Errorf("failed to load config file: %w", err)
		}
		cfg = m.mergeConfigs(cfg, fileCfg)
	}

	return cfg, nil
}

// loadFromFile loads configuration from a TOML file
func (m *ConfigManager) loadFromFile() (Config, error) {
	var cfg Config

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// mergeConfigs merges two config structs, with the second taking precedence
func (m *ConfigManager) mergeConfigs(base, override Config) Config {
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
	if override.Rate != 0 {
		result.Rate = override.Rate
	}

	return result
}

// Save writes the configuration to a TOML file
func (m *ConfigManager) Save(cfg Config) error {
	return m.SaveWithForce(cfg, false)
}

// SaveWithForce writes the configuration to a TOML file
// If force is true, it will overwrite an existing config file
func (m *ConfigManager) SaveWithForce(cfg Config, force bool) error {
	// Check if config file already exists
	if !force {
		if _, err := os.Stat(m.configPath); err == nil {
			return fmt.Errorf("config file %s already exists. Use force to overwrite", m.configPath)
		}
	}

	// Ensure the directory exists
	if dir := filepath.Dir(m.configPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	buffer, err := m.encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	err = os.WriteFile(m.configPath, buffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// InitDefaults creates a config file with default values
func (m *ConfigManager) InitDefaults() error {
	return m.InitDefaultsWithForce(false)
}

// InitDefaultsWithForce creates a config file with default values
// If force is true, it will overwrite an existing config file
func (m *ConfigManager) InitDefaultsWithForce(force bool) error {
	defaults := Defaults()
	return m.SaveWithForce(defaults, force)
}

// encode converts the config into a TOML buffer
func (m *ConfigManager) encode(cfg Config) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	err := toml.NewEncoder(&buffer).Encode(cfg)
	if err != nil {
		return buffer, fmt.Errorf("failed to encode config to TOML: %w", err)
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

	if cfg.Rate <= 0 {
		errs = append(errs, "rate must be greater than 0")
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
		return bytes.Buffer{}, fmt.Errorf("cannot encode config: %w", err)
	}

	return buffer, nil
}

// WriteToToml writes config buffer to file
// Deprecated: Use Manager.Save() instead
func WriteToToml(buffer bytes.Buffer) error {
	return os.WriteFile("config.toml", buffer.Bytes(), 0644)
}

// InitConfigFile creates config with defaults
// Deprecated: Use Manager.InitDefaults() instead
func InitConfigFile() error {
	return InitConfigFileWithForce(false)
}

// InitConfigFileWithForce creates config with defaults and force option
// Deprecated: Use Manager.InitDefaultsWithForce() instead
func InitConfigFileWithForce(force bool) error {
	manager := DefaultConfigManager()
	return manager.InitDefaultsWithForce(force)
}
