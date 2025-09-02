/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com
*/

// Package cmd defines the root command and wires together all
// subcommands for the scrapgo CLI.
package cmd

import (
	"os"

	"github.com/JesterSe7en/scrapgo/cmd/cache"
	"github.com/JesterSe7en/scrapgo/cmd/config"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scrapgo",
	Short: "Scrape one or more URLs concurrently and save results",
	Long: `Scrape web pages or APIs from a list of URLs, using a concurrent worker pool.
Supports rate limiting, retries, and caching of results to avoid redundant requests.
Output can be saved in JSON or CSV format, and verbose logging is available for progress tracking.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Hello from root command")
	// },
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
	rootCmd.AddCommand(cache.NewCommand())
	rootCmd.AddCommand(config.NewCommand())
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scrapgo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
