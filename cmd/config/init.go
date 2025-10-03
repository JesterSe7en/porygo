// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"

	"github.com/JesterSe7en/scrapego/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [filename]",
	Short: "Initialize a default config file with default settings. Defaults to config.toml",
	Long: `The init command creates a new config file in the current directory
with default settings for all available options.

If a filename is provided, it will be used. Otherwise, it defaults to 'config.toml'.

Use this command if you want to generate a fresh configuration file. You can then edit the file
manually or override its values using command-line flags.

Examples:
  scrapego config init
  # creates config.toml with default values

  scrapego config init my-config.toml
  # creates my-config.toml with default values`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := "config.toml"
		if len(args) > 0 {
			configPath = args[0]
		}

		manager := config.NewManager(configPath)

		err := manager.InitDefaults()
		if err != nil {
			fmt.Printf("error creating config file: %v\n", err)
			return
		}

	},
}

func init() {
	configCmd.AddCommand(initCmd)
}
