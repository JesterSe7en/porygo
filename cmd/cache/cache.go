/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com

This command will manage cache results
*/

// Package cache provides the 'cache' command for inspecting,
// clearing, or summarizing cached scraping results.
package cache

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage cached scraping results",
	Long: `Inspect, clear, or summarize previously scraped results stored in the local cache.
This helps avoid unnecessary network requests and enables quick access to past data.
Subcommands include 'list' to view cache entries, 'clear' to remove entries, and 'stats' to view cache statistics.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cache called")
	},
}

// NewCommand returns the cache command for inspecting, clearing,
// or summarizing cached scraping results.
func NewCommand() *cobra.Command {
	return cacheCmd
}

// func init() {
// Here you will define your flags and configuration settings.

// Cobra supports Persistent Flags which will work for this command
// and all subcommands, e.g.:
// cacheCmd.PersistentFlags().String("foo", "", "A help for foo")

// Cobra supports local flags which will only run when this command
// is called directly, e.g.:
// cacheCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
