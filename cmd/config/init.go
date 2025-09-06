/*
Copyright © 2025 Alexander Chan alyxchan87@gmail.com
*/

package config

import (
	"fmt"

	"github.com/JesterSe7en/scrapgo/config"
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
		fmt.Println("Initializing config.toml with default values...")

		force, _ := cmd.Flags().GetBool("force")
		manager := config.DefaultConfigManager()

		err := manager.InitDefaultsWithForce(force)
		if err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			return
		}

		fmt.Println("Successfully created config.toml with default settings")
		fmt.Println("You can now edit the file manually or use command-line flags to override values")
	},
}

func init() {
	configCmd.AddCommand(initCmd)

	// Add force flag to allow overwriting existing config files
	initCmd.Flags().BoolP("force", "f", false, "overwrite existing config.toml file")
}
