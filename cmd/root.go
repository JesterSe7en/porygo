/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com
*/

// Package cmd defines the root command and wires together all
// subcommands for the scrapgo CLI.
package cmd

import (
	"bufio"
	"fmt"
	"os"

	cacheCmd "github.com/JesterSe7en/scrapego/cmd/cache"
	configCmd "github.com/JesterSe7en/scrapego/cmd/config"
	"github.com/JesterSe7en/scrapego/config"
	"github.com/JesterSe7en/scrapego/internal/flags"
	"github.com/JesterSe7en/scrapego/internal/logger"
	"github.com/JesterSe7en/scrapego/internal/scraper"
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
		setupLogging(cmd)
		defer logger.Sync()

		cfg, err := setupConfig(cmd)
		if err != nil {
			return err
		}
		logger.Info("scraping with config : %+v", cfg)

		urls, err := getURLs()
		if err != nil {
			return err
		}

		if len(urls) == 0 {
			logger.Info("no URLs provided via stdin or arguments, nothing to do")
			return nil
		}

		processURLs(cfg, urls)

		return nil
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
	rootCmd.PersistentFlags().String(flags.FlagConfig, "", "specifiy config file")
	rootCmd.Flags().IntP(flags.FlagConcurrency, "c", defaults.Concurrency, "number of workers")
	rootCmd.Flags().DurationP(flags.FlagTimeout, "t", defaults.Timeout, "request timeout per URL")
	rootCmd.Flags().IntP(flags.FlagRetry, "r", defaults.Retry, "number of retries per URL on failure")
	rootCmd.Flags().IntP(flags.FlagBackoff, "b", int(defaults.Backoff), "backoff time between retries")
	rootCmd.Flags().BoolP(flags.FlagForce, "f", defaults.Force, "ignore cache and scrape fresh data")
}

func setupLogging(cmd *cobra.Command) {
	verbose, _ := cmd.PersistentFlags().GetBool(flags.FlagVerbose)
	filename, _ := cmd.PersistentFlags().GetString(flags.FlagLog)
	debug, _ := cmd.PersistentFlags().GetBool(flags.FlagDebug)
	logger.InitLogger(filename, verbose, debug)
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
		logger.Info("using default config values for scraping")
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
	if cmd.Flags().Changed(flags.FlagBackoff) {
		cfg.Backoff, _ = cmd.Flags().GetDuration(flags.FlagBackoff)
	}
	if cmd.Flags().Changed(flags.FlagForce) {
		cfg.Force, _ = cmd.Flags().GetBool(flags.FlagForce)
	}
	return cfg
}

func getURLs() ([]string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %s", err.Error())
	}

	urls := []string{}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading stdin: %v", err)
		}
	}

	return urls, nil
}

func processURLs(cfg config.Config, urls []string) {
	// For now, keep the buffer size same as worker count
	// TODO: evaluate if making the buffer 2x or 3x is worth it
	pool := wp.New(cfg.Concurrency, cfg.Concurrency)
	pool.Run(cfg.Concurrency)

	for _, url := range urls {
		u := url
		pool.Submit(func() wp.Result {
			logger.Info("attempting to scrape: %s", u)
			return scraper.ScrapeWithRetry(u, cfg.Timeout, cfg.Retry, cfg.Backoff)
		})
	}

	pool.Close()

	for res := range pool.Results() {
		if res.Err != nil {
			logger.Error("failed to get response: %s", res.Err.Error())
			continue
		}
		logger.Info("found something: %s", res.Value)
	}
}
