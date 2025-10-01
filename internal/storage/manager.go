// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

package storage

import (
	"sync"
)

// CacheManager manages a single shared cache instance across the application.
// It ensures thread safety and guarantees only one cache instance exists.
type CacheManager struct {
	cache CacheStorage
	mu    sync.RWMutex
}

var (
	manager     *CacheManager
	managerOnce sync.Once
)

// GetCacheManager returns the global cache manager instance.
func GetCacheManager() *CacheManager {
	managerOnce.Do(func() {
		manager = &CacheManager{}
	})
	return manager
}

// GetCache returns the shared cache instance, creating it if needed.
// Thread-safe and ensures only one cache instance exists.
func (m *CacheManager) GetCache() (CacheStorage, error) {
	m.mu.RLock()
	if m.cache != nil {
		defer m.mu.RUnlock()
		return m.cache, nil
	}
	m.mu.RUnlock()

	// Need write lock to create cache
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check in case another goroutine created it
	if m.cache != nil {
		return m.cache, nil
	}

	cache, err := NewBoltCache()
	if err != nil {
		return nil, err
	}

	m.cache = cache
	return m.cache, nil
}

// Close closes the managed cache instance.
func (m *CacheManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cache == nil {
		return nil
	}

	err := m.cache.Close()
	m.cache = nil
	return err
}

// Reset clears the manager state (useful for tests).
func (m *CacheManager) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cache != nil {
		if err := m.cache.Close(); err != nil {
			return err
		}
	}

	m.cache = nil
	return nil
}
