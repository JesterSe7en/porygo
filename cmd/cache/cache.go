/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com

This command will manage cache results
*/

package cache

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cache",
		Short: "Manage cached scraping results",
		Long: `Inspect, clear, or summarize previously scraped results stored in the local cache.
This helps avoid unnecessary network requests and enables quick access to past data.
Subcommands include 'list' to view cache entries, 'clear' to remove entries, and 'stats' to view cache statistics.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cache called")
		},
	}
}

//
// var listCmd = &cobra.Command{
// 	Use:   "list",
// 	Short: "List out the cache",
// 	Long:  `List out the cache long discription goes here`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("You are asking to list the cache")
// 	},
// }
//
// var clearCmd = &cobra.Command{
// 	Use:   "clear",
// 	Short: "Clear all or by URL",
// 	Long:  `Clear entire cache or by specific URL`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("You are asking to list the cache")
// 	},
// }
//
// var statsCmd = &cobra.Command{
// 	Use:   "stats",
// 	Short: "Summary of cache (entries, last updated, etc.)",
// 	Long:  `Show the entire summary of the cache`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("You are asking to list the cache")
// 	},
// }

// func init() {
// Here you will define your flags and configuration settings.

// Cobra supports Persistent Flags which will work for this command
// and all subcommands, e.g.:
// cacheCmd.PersistentFlags().String("foo", "", "A help for foo")

// Cobra supports local flags which will only run when this command
// is called directly, e.g.:
// cacheCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
