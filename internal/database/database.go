// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

// Package database is a wrapper for bbolt to provide a simple interface to the database
package database

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"go.etcd.io/bbolt"
)

var (
	db         *bbolt.DB
	dbOnce     sync.Once
	dbPath     string
	bucketName = []byte("cache")
)

// getCachePath determines the cache path cross-platform.
func getCachePath() (string, error) {
	if x := os.Getenv("XDG_CACHE_HOME"); x != "" {
		return filepath.Join(x, "scrapego", "cache.db"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Local", "Scrapego", "cache.db"), nil
	default: // unix-like
		return filepath.Join(home, ".cache", "scrapego", "cache.db"), nil
	}
}

// getDB lazily opens the DB if it hasn’t been opened yet.
func getDB() (*bbolt.DB, error) {
	var err error
	// sync.Once to only initialize the db if it hasn't already
	dbOnce.Do(func() {
		dbPath, err = getCachePath()
		if err != nil {
			return
		}

		if err = os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return
		}

		db, err = bbolt.Open(dbPath, 0o600, nil)
		if err != nil {
			return
		}

		// Ensure bucket exists
		err = db.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucketName)
			return err
		})
	})
	return db, err
}

// Put stores a key-value pair.
func Put(key, value []byte) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket does not exist")
		}
		return b.Put(key, value)
	})
}

// Get retrieves a value by key.
func Get(key []byte) ([]byte, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

	var val []byte
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket does not exist")
		}
		val = b.Get(key)
		return nil
	})
	return val, err
}

// Delete removes a key from the bucket.
func Delete(key []byte) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("bucket does not exist")
		}
		return b.Delete(key)
	})
}

func ClearCache() error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

// Close closes the database.
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
