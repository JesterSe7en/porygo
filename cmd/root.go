/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com
*/

// Package cmd defines the root command and wires together all
// subcommands for the scrapgo CLI.
package cmd

import (
	"fmt"
	"os"

	cacheCmd "github.com/JesterSe7en/scrapgo/cmd/cache"
	configCmd "github.com/JesterSe7en/scrapgo/cmd/config"
	"github.com/JesterSe7en/scrapgo/config"
	l "github.com/JesterSe7en/scrapgo/internal/logger"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scrapgo",
	Short: "Scrape one or more URLs concurrently and save results",
	Long: `Scrape web pages or APIs from a list of URLs, using a concurrent worker pool.
Supports rate limiting, retries, and caching of results to avoid redundant requests.
Output can be saved in JSON or CSV format, and verbose logging is available for progress tracking.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration with proper precedence
		manager := config.DefaultConfigManager()
		cfg, err := manager.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Override with CLI flags if provided
		cfg = mergeCLIFlags(cmd, cfg)

		// Validate configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		// TODO: Implement actual scraping logic here
		fmt.Printf("Would scrape with config: %+v\n", cfg)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	l.Debug("Executing root command...")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// mergeCLIFlags merges CLI flag values into the configuration
func mergeCLIFlags(cmd *cobra.Command, cfg config.Config) config.Config {
	if cmd.PersistentFlags().Changed("verbose") {
		cfg.Verbose, _ = cmd.Flags().GetBool("verbose")
	}
	if cmd.Flags().Changed("input") {
		cfg.Input, _ = cmd.Flags().GetString("input")
	}
	if cmd.Flags().Changed("concurrency") {
		cfg.Concurrency, _ = cmd.Flags().GetInt("concurrency")
	}
	if cmd.Flags().Changed("timeout") {
		cfg.Timeout, _ = cmd.Flags().GetDuration("timeout")
	}
	if cmd.Flags().Changed("output") {
		cfg.Output, _ = cmd.Flags().GetString("output")
	}
	if cmd.Flags().Changed("retry") {
		cfg.Retry, _ = cmd.Flags().GetInt("retry")
	}
	if cmd.Flags().Changed("rate") {
		cfg.Rate, _ = cmd.Flags().GetInt("rate")
	}
	if cmd.Flags().Changed("force") {
		cfg.Force, _ = cmd.Flags().GetBool("force")
	}
	return cfg
}

func init() {
	rootCmd.AddCommand(cacheCmd.NewCommand())
	rootCmd.AddCommand(configCmd.NewCommand())

	// Get default values for flag defaults
	defaults := config.Defaults()

	// Define flags with default values
	rootCmd.PersistentFlags().BoolP("verbose", "v", defaults.Verbose, "show logs for each step")
	rootCmd.Flags().StringP("input", "i", defaults.Input, "path to file with URLs")
	rootCmd.Flags().IntP("concurrency", "c", defaults.Concurrency, "number of workers")
	rootCmd.Flags().DurationP("timeout", "t", defaults.Timeout, "request timeout per URL")
	rootCmd.Flags().StringP("output", "o", defaults.Output, "JSON or CSV")
	rootCmd.Flags().IntP("retry", "r", defaults.Retry, "number of retries per URL on failure")
	rootCmd.Flags().IntP("rate", "R", defaults.Rate, "requests per second")
	rootCmd.Flags().BoolP("force", "f", defaults.Force, "ignore cache and scrape fresh data")
}
