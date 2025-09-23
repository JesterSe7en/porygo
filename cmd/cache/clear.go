// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

package cache

import (
	"github.com/JesterSe7en/scrapego/internal/database"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the local cache of all scraped data.",
	Long: `The clear command removes all stored responses from the local cache.

By default, scrapego caches successful responses to avoid re-fetching the same URLs.
Running this command will empty the cache, ensuring that the next scrape for any URL
will fetch fresh data from the source. This is useful if the content of the URLs
has changed and you need to force an update.

Example:
  scrapego cache clear`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.ClearCache()
	},
}

func init() {
	cacheCmd.AddCommand(clearCmd)
}
