// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

package cache

import (
	"context"
	"fmt"

	"github.com/JesterSe7en/porygo/internal/storage"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the local cache of all scraped data.",
	Long: `The clear command removes all stored responses from the local cache.

By default, porygo caches successful responses to avoid re-fetching the same URLs.
Running this command will empty the cache, ensuring that the next scrape for any URL
will fetch fresh data from the source. This is useful if the content of the URLs
has changed and you need to force an update.

Example:
  porygo cache clear`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := storage.GetCacheManager()
		cache, err := manager.GetCache()
		if err != nil {
			return fmt.Errorf("failed to get cache: %w", err)
		}

		if err := cache.Clear(context.Background()); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}

		fmt.Println("Cache cleared successfully.")
		return nil
	},
}
