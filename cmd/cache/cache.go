// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

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
	Long: `This command provides tools for clearing cached scraping results.
This helps avoid unnecessary network requests and enables quick access to past data.
Subcommands include 'clear' to remove entries.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cache called")
	},
}

// NewCommand returns the cache command for inspecting, clearing,
// or summarizing cached scraping results.
func NewCommand() *cobra.Command {
	cacheCmd.AddCommand(clearCmd)
	return cacheCmd
}
