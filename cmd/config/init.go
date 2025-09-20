// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"

	"github.com/JesterSe7en/scrapego/config"
	"github.com/JesterSe7en/scrapego/internal/logger"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a default config.toml file with default settings",
	Long: `The init command creates a new config.toml file in the current directory
with default settings for all available options.

Use this command if you want to generate a fresh configuration file. You can then edit the file
manually or override its values using command-line flags.

Example:
  scrapego config init
  # creates config.toml with default values`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("initializing config.toml with default values")

		force, _ := cmd.Flags().GetBool("force")
		manager := config.DefaultManager()

		err := manager.InitDefaultsWithForce(force)
		if err != nil {
			fmt.Printf("error creating config file: %v\n", err)
			return
		}

		logger.Info("successfully created config.toml with default settings")
	},
}

func init() {
	configCmd.AddCommand(initCmd)

	// Add force flag to allow overwriting existing config files
	initCmd.Flags().BoolP("force", "f", false, "overwrite existing config.toml file")
}
