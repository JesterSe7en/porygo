// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

package cache

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cache/statsCmd represents the cache/stats command
var statsCmd = &cobra.Command{
	Use:   "cache/stats",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cache/stats called")
	},
}

func init() {
	cacheCmd.AddCommand(statsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cache/statsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cache/statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
