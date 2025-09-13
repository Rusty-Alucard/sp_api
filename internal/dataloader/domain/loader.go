// Package dataloader provides DataLoader domain models and interfaces.
package dataloader

import (
	"context"
	"time"
)

const (
	// DefaultMaxBatch is the default maximum batch size
	DefaultMaxBatch = 100
	// DefaultWait is the default wait time for batching
	DefaultWait = 4 * time.Millisecond
	// DefaultCacheExpiry is the default cache expiry time
	DefaultCacheExpiry = 5 * time.Minute
)

// Loader defines the interface for DataLoader operations.
type Loader[K comparable, V any] interface {
	Load(ctx context.Context, key K) (V, error)
	LoadMany(ctx context.Context, keys []K) ([]V, error)
	Clear(key K)
	ClearAll()
	Prime(key K, value V)
}

// BatchFunc defines the function signature for batch loading.
type BatchFunc[K comparable, V any] func(ctx context.Context, keys []K) ([]V, error)

// Config represents DataLoader configuration.
type Config struct {
	MaxBatch    int
	Wait        time.Duration
	Cache       bool
	CacheExpiry time.Duration
}

// DefaultConfig returns the default DataLoader configuration.
func DefaultConfig() Config {
	return Config{
		MaxBatch:    DefaultMaxBatch,
		Wait:        DefaultWait,
		Cache:       true,
		CacheExpiry: DefaultCacheExpiry,
	}
}

// Cache represents a cache interface for DataLoader.
type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V, expiry time.Duration)
	Delete(key K)
	Clear()
}

// MemoryCache implements an in-memory cache for DataLoader.
type MemoryCache[K comparable, V any] struct {
	data map[K]cacheEntry[V]
}

type cacheEntry[V any] struct {
	value     V
	expiresAt time.Time
}

// NewMemoryCache creates a new memory cache.
func NewMemoryCache[K comparable, V any]() *MemoryCache[K, V] {
	return &MemoryCache[K, V]{
		data: make(map[K]cacheEntry[V]),
	}
}

// Get retrieves a value from the cache.
func (c *MemoryCache[K, V]) Get(key K) (V, bool) {
	entry, exists := c.data[key]
	if !exists {
		var zero V
		return zero, false
	}

	if time.Now().After(entry.expiresAt) {
		delete(c.data, key)
		var zero V
		return zero, false
	}

	return entry.value, true
}

// Set stores a value in the cache.
func (c *MemoryCache[K, V]) Set(key K, value V, expiry time.Duration) {
	c.data[key] = cacheEntry[V]{
		value:     value,
		expiresAt: time.Now().Add(expiry),
	}
}

// Delete removes a value from the cache.
func (c *MemoryCache[K, V]) Delete(key K) {
	delete(c.data, key)
}

// Clear removes all values from the cache.
func (c *MemoryCache[K, V]) Clear() {
	c.data = make(map[K]cacheEntry[V])
}

