/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com

View and modify the CLI configuration
*/

package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "View and modify CLI configuration",
		Long: `View or update the scraper's configuration settings, such as default concurrency,
rate limits, output paths, or user-agent strings.
Supports a config file (YAML) to persist settings across sessions.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("config called")
		},
	}
}

// func init() {
// 	// Here you will define your flags and configuration settings.

// Cobra supports Persistent Flags which will work for this command
// and all subcommands, e.g.:
// configCmd.PersistentFlags().String("foo", "", "A help for foo")

// Cobra supports local flags which will only run when this command
// is called directly, e.g.:
// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
