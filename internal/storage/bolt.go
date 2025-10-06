package storage

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.etcd.io/bbolt"
	bboltErrors "go.etcd.io/bbolt/errors"
)

const (
	// Default cache file permissions
	cacheFileMode = 0o600

	// Default bucket name for cache entries
	defaultBucketName = "cache"
)

// Default bucket name as byte slice since bbolt requires byte slices for bucket names
var bucketName = []byte(defaultBucketName)

// Common errors
var (
	ErrNotFound       = errors.New("entry not found")
	ErrBucketNotFound = errors.New("bucket not found")
	ErrEncoding       = errors.New("failed to encode cache entry")
	ErrDecoding       = errors.New("failed to decode cache entry")
)

type boltCache struct {
	db *bbolt.DB
}

// getCachePath determines the appropriate cache directory path for the current platform.
// It follows XDG Base Directory specification on Unix-like systems and uses appropriate
// directories on Windows.
func getCachePath() (string, error) {
	// Check for XDG_CACHE_HOME environment variable first
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
		return filepath.Join(xdgCache, "porygo", "cache.db"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Local", "porygo", "cache.db"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Caches", "porygo", "cache.db"), nil
	default: // Unix-like systems
		return filepath.Join(home, ".cache", "porygo", "cache.db"), nil
	}
}

func newBoltCacheAt(pathDB string) (CacheStorage, error) {
	// Ensure the directory exists before opening the database.
	if err := os.MkdirAll(filepath.Dir(pathDB), 0o750); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	db, err := bbolt.Open(pathDB, cacheFileMode, &bbolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB at %s: %w", pathDB, err)
	}

	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
		return nil
	}); err != nil {
		db.Close()
		return nil, err
	}

	return &boltCache{db: db}, nil
}

func NewBoltCache() (CacheStorage, error) {
	pathDB, err := getCachePath()
	if err != nil {
		return nil, fmt.Errorf("cannot get cache location: %w", err)
	}
	return newBoltCacheAt(pathDB)
}

// Get retrieves a cache entry by key.
// Returns ErrNotFound if the key doesn't exist.
func (b *boltCache) Get(ctx context.Context, key string) (CacheEntry, error) {
	if len(key) == 0 {
		return CacheEntry{}, errors.New("key cannot be empty")
	}

	var value []byte

	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		value = bucket.Get([]byte(key))
		return nil
	})
	if err != nil {
		return CacheEntry{}, fmt.Errorf("failed to read from database: %w", err)
	}

	if value == nil {
		return CacheEntry{}, ErrNotFound
	}

	entry, err := decodeEntry(value)
	if err != nil {
		return CacheEntry{}, fmt.Errorf("failed to decode entry: %w", err)
	}

	return entry, nil
}

// Put stores a key-value pair in the cache with a current timestamp.
func (b *boltCache) Set(ctx context.Context, key string, entry CacheEntry) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	encodedEntry, err := encodeEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to encode entry: %w", err)
	}

	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		if err := bucket.Put([]byte(key), encodedEntry); err != nil {
			return fmt.Errorf("failed to put key-value pair: %w", err)
		}

		return nil
	})
}

// Delete removes a key from the cache.
func (b *boltCache) Delete(ctx context.Context, key string) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		if err := bucket.Delete([]byte(key)); err != nil {
			return fmt.Errorf("failed to delete key: %w", err)
		}

		return nil
	})
}

// ClearCache removes all entries from the cache by recreating the bucket.
func (b *boltCache) Clear(ctx context.Context) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		// Delete existing bucket
		if err := tx.DeleteBucket(bucketName); err != nil {
			// Ignore error if bucket doesn't exist
			if !errors.Is(err, bboltErrors.ErrBucketNotFound) {
				return fmt.Errorf("failed to delete bucket: %w", err)
			}
		}

		// Recreate bucket
		if _, err := tx.CreateBucket(bucketName); err != nil {
			return fmt.Errorf("failed to recreate bucket: %w", err)
		}

		return nil
	})
}

// Close closes the database connection.
func (b *boltCache) Close() error {
	if b.db != nil {
		err := b.db.Close()
		b.db = nil
		return err
	}
	return nil
}

// encodeEntry serializes a CacheEntry using gob encoding.
func encodeEntry(entry CacheEntry) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(entry); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrEncoding, err)
	}

	return buf.Bytes(), nil
}

// decodeEntry deserializes a CacheEntry from gob-encoded data.
func decodeEntry(data []byte) (CacheEntry, error) {
	var entry CacheEntry
	decoder := gob.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(&entry); err != nil {
		return CacheEntry{}, fmt.Errorf("%w: %w", ErrDecoding, err)
	}

	return entry, nil
}
