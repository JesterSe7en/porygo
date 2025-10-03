// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package cmd defines the root command and wires together all
// subcommands for the scrapgo CLI.
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"

	cacheCmd "github.com/JesterSe7en/scrapego/cmd/cache"
	configCmd "github.com/JesterSe7en/scrapego/cmd/config"
	"github.com/JesterSe7en/scrapego/config"

	"github.com/JesterSe7en/scrapego/internal/app"
	"github.com/JesterSe7en/scrapego/internal/flags"
	"github.com/JesterSe7en/scrapego/internal/logger"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scrapego [urls...]",
	Short: "Scrape one or more URLs concurrently and save results",
	Long: `Scrape web pages or APIs from a list of URLs, using a concurrent worker pool.
Supports rate limiting, retries, and caching of results to avoid redundant requests.
Output can be saved in JSON or CSV format, and verbose logging is available for progress tracking.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// RunE will only grab flags and parse them into config; this includes list of URLs
		verbose, _ := cmd.PersistentFlags().GetBool(flags.FlagVerbose)
		filename, _ := cmd.PersistentFlags().GetString(flags.FlagLog)
		debug, _ := cmd.PersistentFlags().GetBool(flags.FlagDebug)

		log, err := logger.New(filename, debug, verbose)
		if err != nil {
			return err
		}
		defer log.Sync()

		cfg, err := setupConfig(cmd)
		log.Debug("scraping with config : %+v", cfg)
		if err != nil {
			return err
		}

		urls, err := getURLs(args)
		if err != nil {
			return err
		}

		if len(urls) == 0 {
			return cmd.Help()
		}

		app, err := app.New(&log, &cfg)
		if err != nil {
			return err
		}

		return app.Run(context.Background(), urls)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cacheCmd.NewCommand())
	rootCmd.AddCommand(configCmd.NewCommand())

	// Get default values for flag defaults
	defaults := config.Defaults()

	// Define flags with default values
	// log, debug, and verbose is not in the Defaults struct as that is used to init a config.toml file
	// do not wnat those to be exposed in config.  user wil have to specifiy these flags explicity during the command call
	rootCmd.PersistentFlags().StringP(flags.FlagLog, "l", "", "file path to write logs")
	rootCmd.PersistentFlags().BoolP(flags.FlagDebug, "d", false, "output debug messages")
	rootCmd.PersistentFlags().BoolP(flags.FlagVerbose, "v", false, "show logs for each step")
	// config and Concurrency cannot use same shorthand character
	rootCmd.PersistentFlags().String(flags.FlagConfig, "", "specify config file")
	rootCmd.Flags().IntP(flags.FlagConcurrency, "c", defaults.Concurrency, "number of workers")
	rootCmd.Flags().DurationP(flags.FlagTimeout, "t", defaults.Timeout, "request timeout per URL")
	rootCmd.Flags().IntP(flags.FlagRetry, "r", defaults.Retry, "number of retries per URL on failure")
	rootCmd.Flags().Duration(flags.FlagRetryDelay, defaults.Backoff.BaseDelay, "base delay between retries (exponential backoff applied)")
	rootCmd.Flags().Bool(flags.FlagRetryJitter, defaults.Backoff.Jitter, "enable jitter for retry delays")
	rootCmd.Flags().BoolP(flags.FlagForce, "f", defaults.Force, "ignore cache and scrape fresh data")

	// scraper flags
	rootCmd.Flags().StringSliceP(flags.FlagSelect, "s", []string{}, "CSS selectors to extract")
	rootCmd.Flags().StringSliceP(flags.FlagPattern, "p", []string{}, "regex patterns to match")
	rootCmd.Flags().StringP(flags.FlagFormat, "o", "json", "output format (json|csv|plain)")
	rootCmd.Flags().BoolP(flags.FlagQuiet, "q", false, "only output extracted data")
	rootCmd.Flags().BoolP(flags.FlagHeaders, "H", false, "include response headers")

}

func setupConfig(cmd *cobra.Command) (config.Config, error) {
	manager := config.DefaultManager()

	configFile, err := cmd.Flags().GetString(flags.FlagConfig)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to load configuration: %s", err.Error())
	}

	var cfg config.Config
	// Load configuration with proper precedence
	if configFile == "" {
		cfg = manager.LoadDefaults()
	} else {
		cfg, err = manager.LoadFromFile(configFile)
		if err != nil {
			return config.Config{}, fmt.Errorf("failed to load configuration: %s", err.Error())
		}
	}

	// Flags manually set takes precedence over whatever config file says
	// Override with CLI flags if provided
	cfg = mergeCLIFlags(cmd, cfg)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return config.Config{}, fmt.Errorf("invalid configuration: %s", err.Error())
	}

	return cfg, nil
}

// mergeCLIFlags merges CLI flag values into the configuration
func mergeCLIFlags(cmd *cobra.Command, cfg config.Config) config.Config {
	if cmd.Flags().Changed(flags.FlagConcurrency) {
		cfg.Concurrency, _ = cmd.Flags().GetInt(flags.FlagConcurrency)
	}
	if cmd.Flags().Changed(flags.FlagTimeout) {
		cfg.Timeout, _ = cmd.Flags().GetDuration(flags.FlagTimeout)
	}
	if cmd.Flags().Changed(flags.FlagRetry) {
		cfg.Retry, _ = cmd.Flags().GetInt(flags.FlagRetry)
	}
	if cmd.Flags().Changed(flags.FlagRetryDelay) {
		cfg.Backoff.BaseDelay, _ = cmd.Flags().GetDuration(flags.FlagRetryDelay)
	}
	if cmd.Flags().Changed(flags.FlagRetryJitter) {
		cfg.Backoff.Jitter, _ = cmd.Flags().GetBool(flags.FlagRetryJitter)
	}
	if cmd.Flags().Changed(flags.FlagForce) {
		cfg.Force, _ = cmd.Flags().GetBool(flags.FlagForce)
	}

	// scraper flags
	if cmd.Flags().Changed(flags.FlagSelect) {
		cfg.SelectorsConfig.Select, _ = cmd.Flags().GetStringSlice(flags.FlagSelect)
	}
	if cmd.Flags().Changed(flags.FlagPattern) {
		cfg.SelectorsConfig.Pattern, _ = cmd.Flags().GetStringSlice(flags.FlagPattern)
	}
	if cmd.Flags().Changed(flags.FlagFormat) {
		cfg.Format, _ = cmd.Flags().GetString(flags.FlagFormat)
	}
	if cmd.Flags().Changed(flags.FlagQuiet) {
		cfg.Quiet, _ = cmd.Flags().GetBool(flags.FlagQuiet)
	}
	if cmd.Flags().Changed(flags.FlagHeaders) {
		cfg.Headers, _ = cmd.Flags().GetBool(flags.FlagHeaders)
	}

	return cfg
}

func getURLs(args []string) ([]string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %s", err.Error())
	}

	// Check for stdin first
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		urls := []string{}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading stdin: %v", err)
		}
		if len(urls) > 0 {
			err := validateURLs(urls)
			if err != nil {
				return nil, err
			}
			return urls, nil
		}
	}

	// If no stdin, use args
	if len(args) > 0 {
		err := validateURLs(args)
		if err != nil {
			return nil, err
		}
		return args, nil
	}

	// No input from stdin or args
	return []string{}, nil
}

func validateURLs(inputs []string) error {
	for _, input := range inputs {
		if _, err := url.Parse(input); err != nil {
			return fmt.Errorf("invalid URL: %s", input)
		}
	}
	return nil
}
