// Package storage contains the interface that is used to store and retrieve cache entries
package storage

import (
	"context"
	"time"
)

type CacheEntry struct {
	Value          []byte
	ExpirationTime time.Time
}

type CacheStorage interface {
	Get(ctx context.Context, key string) (CacheEntry, error)
	Set(ctx context.Context, key string, value CacheEntry) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Close() error
}
