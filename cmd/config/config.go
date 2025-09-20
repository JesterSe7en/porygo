// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package config provides the 'config' command for viewing and
// modifying the CLI configuration settings.
package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View and modify CLI configuration",
	Long: `View or update the scraper's configuration settings, such as default concurrency,
rate limits, output paths, or user-agent strings.
Supports a config file (YAML) to persist settings across sessions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

// NewCommand returns the config command for viewing and modifying
// the CLI configuration.
func NewCommand() *cobra.Command {
	return configCmd
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
