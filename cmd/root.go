/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com
*/

// Package cmd defines the root command and wires together all
// subcommands for the scrapgo CLI.
package cmd

import (
	"fmt"
	"os"
	"time"

	cacheCmd "github.com/JesterSe7en/scrapgo/cmd/cache"
	configCmd "github.com/JesterSe7en/scrapgo/cmd/config"
	"github.com/JesterSe7en/scrapgo/config"
	"github.com/spf13/cobra"
)

var cfg config.Config

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scrapgo",
	Short: "Scrape one or more URLs concurrently and save results",
	Long: `Scrape web pages or APIs from a list of URLs, using a concurrent worker pool.
Supports rate limiting, retries, and caching of results to avoid redundant requests.
Output can be saved in JSON or CSV format, and verbose logging is available for progress tracking.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Config: %+v\n", cfg)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scrapgo.yaml)")
	// flags
	//
	// --input, -i → path to file with URLs
	// 	•	--concurrency, -c → number of workers (default 5)
	// 	•	--timeout, -t → request timeout per URL (default 10s)
	// 	•	--output, -o → JSON or CSV (default JSON)
	// 	•	--verbose, -v → show logs for each step
	// 	•	--retry, -r → number of retries on failure (default 3)
	// 	•	--rate, -R → requests per second (default 1)
	// 	•	--force, -f → ignore cache and scrape fresh
	//
	rootCmd.Flags().StringVarP(&cfg.Input, "input", "i", "", "path to file with URLs")
	rootCmd.Flags().IntVarP(&cfg.Concurrency, "concurrency", "c", 5, "number of workers")
	rootCmd.Flags().DurationVarP(&cfg.Timeout, "timeout", "t", 10*time.Second, "requirest timeout per URL")
	rootCmd.Flags().StringVarP(&cfg.Output, "output", "o", "JSON", "JSON or CSV")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "shows logs for each step")
	rootCmd.Flags().IntVarP(&cfg.Retry, "retry", "r", 3, "number of retries per URL on failure")
	rootCmd.Flags().IntVarP(&cfg.Rate, "rate", "R", 1, "requests per second (default 1)")
	rootCmd.Flags().BoolVarP(&cfg.Force, "force", "f", false, "ignore cache and scrape fresh data")
}
