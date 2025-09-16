/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com
*/

// Package cmd defines the root command and wires together all
// subcommands for the scrapgo CLI.
package cmd

import (
	"fmt"
	"os"

	cacheCmd "github.com/JesterSe7en/scrapego/cmd/cache"
	configCmd "github.com/JesterSe7en/scrapego/cmd/config"
	c "github.com/JesterSe7en/scrapego/config"
	f "github.com/JesterSe7en/scrapego/internal/flags"
	l "github.com/JesterSe7en/scrapego/internal/logger"
	s "github.com/JesterSe7en/scrapego/internal/scraper"
	wp "github.com/JesterSe7en/scrapego/internal/workerpool"
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
		manager := c.DefaultManager()
		cfg, err := manager.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %s", err.Error())
		}

		// Override with CLI flags if provided
		cfg = mergeCLIFlags(cmd, cfg)

		// Validate configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %s", err.Error())
		}

		l.Info("Scraping with config : %+v", cfg)

		// For now, keep the buffer size same as worker count
		// TODO: evaluate if making the buffer 2x or 3x is worth it

		pool := wp.New(cfg.Concurrency, cfg.Concurrency)
		pool.Run(cfg.Concurrency)

		pool.Submit(func() wp.Result {
			url := "http://www.google.com"
			l.Info("attempting to scrape: %s", url)
			return s.ScrapeWithRetry(url, cfg.Timeout, cfg.Retry, cfg.Backoff)
		})

		pool.Close()

		for res := range pool.Results() {
			if res.Err != nil {
				l.Error("failed to get response: %s", err.Error())
				continue
			}

			l.Info("found something: %s", res.Value)

		}
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
func mergeCLIFlags(cmd *cobra.Command, cfg c.Config) c.Config {
	if cmd.PersistentFlags().Changed(f.FlagVerbose) {
		cfg.Verbose, _ = cmd.Flags().GetBool(f.FlagVerbose)
	}
	if cmd.Flags().Changed(f.FlagInput) {
		cfg.Input, _ = cmd.Flags().GetString(f.FlagInput)
	}
	if cmd.Flags().Changed(f.FlagConcurrency) {
		cfg.Concurrency, _ = cmd.Flags().GetInt(f.FlagConcurrency)
	}
	if cmd.Flags().Changed(f.FlagTimeout) {
		cfg.Timeout, _ = cmd.Flags().GetDuration(f.FlagTimeout)
	}
	if cmd.Flags().Changed(f.FlagOutput) {
		cfg.Output, _ = cmd.Flags().GetString(f.FlagOutput)
	}
	if cmd.Flags().Changed(f.FlagRetry) {
		cfg.Retry, _ = cmd.Flags().GetInt(f.FlagRetry)
	}
	if cmd.Flags().Changed(f.FlagBackoff) {
		cfg.Backoff, _ = cmd.Flags().GetDuration(f.FlagBackoff)
	}
	if cmd.Flags().Changed(f.FlagForce) {
		cfg.Force, _ = cmd.Flags().GetBool(f.FlagForce)
	}
	return cfg
}

func init() {
	rootCmd.AddCommand(cacheCmd.NewCommand())
	rootCmd.AddCommand(configCmd.NewCommand())

	// Get default values for flag defaults
	defaults := c.Defaults()

	// Define flags with default values
	rootCmd.PersistentFlags().BoolP(f.FlagVerbose, "v", defaults.Verbose, "show logs for each step")
	rootCmd.Flags().StringP(f.FlagInput, "i", defaults.Input, "path to file with URLs")
	rootCmd.Flags().IntP(f.FlagConcurrency, "c", defaults.Concurrency, "number of workers")
	rootCmd.Flags().DurationP(f.FlagTimeout, "t", defaults.Timeout, "request timeout per URL")
	rootCmd.Flags().StringP(f.FlagOutput, "o", defaults.Output, "JSON or CSV")
	rootCmd.Flags().IntP(f.FlagRetry, "r", defaults.Retry, "number of retries per URL on failure")
	rootCmd.Flags().IntP(f.FlagBackoff, "b", int(defaults.Backoff), "backoff time between retries")
	rootCmd.Flags().BoolP(f.FlagForce, "f", defaults.Force, "ignore cache and scrape fresh data")
}
