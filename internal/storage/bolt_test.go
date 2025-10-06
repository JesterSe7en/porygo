package storage

import (
	"bytes"
	"context"
	"path"
	"testing"
	"time"
)

func Test_boltCache(t *testing.T) {
	t.Run("Test NewBoltCache", func(t *testing.T) {
		cache, err := newBoltCacheAt(path.Join(t.TempDir(), "cache.db"))
		if err != nil {
			t.Fatalf("Failed to create new bolt cache: %v", err)
		}
		defer cache.Close()

		if cache == nil {
			t.Fatal("Expected cache to not be nil")
		}
	})

	t.Run("Test Get and Set", func(t *testing.T) {
		cache, err := newBoltCacheAt(path.Join(t.TempDir(), "cache.db"))

		if err != nil {
			t.Fatalf("Failed to create new bolt cache: %v", err)
		}
		defer cache.Close()

		key := "test-key"
		entry := CacheEntry{
			Value:          []byte("test-value"),
			ExpirationTime: time.Now().Add(1 * time.Hour),
		}

		if err := cache.Set(context.Background(), key, entry); err != nil {
			t.Fatalf("Failed to set cache entry: %v", err)
		}

		retrievedEntry, err := cache.Get(context.Background(), key)
		if err != nil {
			t.Fatalf("Failed to get cache entry: %v", err)
		}

		if string(retrievedEntry.Value) != string(entry.Value) {
			t.Errorf("Expected value %s, but got %s", string(entry.Value), string(retrievedEntry.Value))
		}
	})

	t.Run("Test Delete", func(t *testing.T) {
		cache, err := newBoltCacheAt(path.Join(t.TempDir(), "cache.db"))
		if err != nil {
			t.Fatalf("Failed to create new bolt cache: %v", err)
		}
		defer cache.Close()

		key := "test-key"
		entry := CacheEntry{
			Value:          []byte("test-value"),
			ExpirationTime: time.Now().Add(1 * time.Hour),
		}

		if err := cache.Set(context.Background(), key, entry); err != nil {
			t.Fatalf("Failed to set cache entry: %v", err)
		}

		if err := cache.Delete(context.Background(), key); err != nil {
			t.Fatalf("Failed to delete cache entry: %v", err)
		}

		_, err = cache.Get(context.Background(), key)
		if err != ErrNotFound {
			t.Errorf("Expected error %v, but got %v", ErrNotFound, err)
		}
	})

	t.Run("Test Clear", func(t *testing.T) {
		cache, err := newBoltCacheAt(path.Join(t.TempDir(), "cache.db"))
		if err != nil {
			t.Fatalf("Failed to create new bolt cache: %v", err)
		}
		defer cache.Close()

		key1 := "test-key-1"
		entry1 := CacheEntry{
			Value:          []byte("test-value-1"),
			ExpirationTime: time.Now().Add(1 * time.Hour),
		}

		key2 := "test-key-2"
		entry2 := CacheEntry{
			Value:          []byte("test-value-2"),
			ExpirationTime: time.Now().Add(1 * time.Hour),
		}

		if err := cache.Set(context.Background(), key1, entry1); err != nil {
			t.Fatalf("Failed to set cache entry: %v", err)
		}

		if err := cache.Set(context.Background(), key2, entry2); err != nil {
			t.Fatalf("Failed to set cache entry: %v", err)
		}

		if err := cache.Clear(context.Background()); err != nil {
			t.Fatalf("Failed to clear cache: %v", err)
		}

		_, err = cache.Get(context.Background(), key1)
		if err != ErrNotFound {
			t.Errorf("Expected error %v, but got %v", ErrNotFound, err)
		}

		_, err = cache.Get(context.Background(), key2)
		if err != ErrNotFound {
			t.Errorf("Expected error %v, but got %v", ErrNotFound, err)
		}
	})

	t.Run("Test Encode and Decode", func(t *testing.T) {
		entry := CacheEntry{
			Value:          []byte("test-value"),
			ExpirationTime: time.Now().Add(1 * time.Hour),
		}

		data, err := encodeEntry(entry)
		if err != nil {
			t.Fatalf("unexpected error encoding entry: %v", err)
		}

		if len(data) == 0 {
			t.Fatal("expected non-empty encoded data")
		}

		decodedEntry, err := decodeEntry(data)
		if err != nil {
			t.Fatalf("unexpected error decoding entry: %v", err)
		}

		// verify the integrity of the decoded entry
		if !bytes.Equal(entry.Value, decodedEntry.Value) {
			t.Errorf("Expected value %v, but got %v", entry.Value, decodedEntry.Value)
		}

		if !entry.ExpirationTime.Equal(decodedEntry.ExpirationTime) {
			t.Errorf("Expected expiration time %v, but got %v", entry.ExpirationTime, decodedEntry.ExpirationTime)
		}

	})
}
