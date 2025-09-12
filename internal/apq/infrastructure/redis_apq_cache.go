// Package apq provides APQ infrastructure implementations.
package apq

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisAPQCache implements APQ cache using Redis.
type RedisAPQCache struct {
	client redis.UniversalClient
	prefix string
}

// NewRedisAPQCache creates a new RedisAPQCache.
func NewRedisAPQCache(client redis.UniversalClient, prefix string) *RedisAPQCache {
	if prefix == "" {
		prefix = "apq:"
	}

	return &RedisAPQCache{
		client: client,
		prefix: prefix,
	}
}

// Get retrieves a persisted query from the cache.
func (c *RedisAPQCache) Get(ctx context.Context, hash string) (string, bool, error) {
	key := c.prefix + hash

	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return result, true, nil
}

// Set stores a persisted query in the cache.
func (c *RedisAPQCache) Set(ctx context.Context, hash, query string, ttl time.Duration) error {
	key := c.prefix + hash

	return c.client.Set(ctx, key, query, ttl).Err()
}

// Delete removes a persisted query from the cache.
func (c *RedisAPQCache) Delete(ctx context.Context, hash string) error {
	key := c.prefix + hash

	return c.client.Del(ctx, key).Err()
}

// Clear removes all persisted queries from the cache.
func (c *RedisAPQCache) Clear(ctx context.Context) error {
	pattern := c.prefix + "*"

	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.client.Del(ctx, keys...).Err()
}

// GetStats returns cache statistics.
func (c *RedisAPQCache) GetStats(ctx context.Context) (*APQStats, error) {
	pattern := c.prefix + "*"

	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	stats := &APQStats{
		TotalQueries: len(keys),
		CacheSize:    len(keys),
	}

	// Get TTL for each key to calculate average
	totalTTL := time.Duration(0)
	validKeys := 0

	for _, key := range keys {
		ttl, err := c.client.TTL(ctx, key).Result()
		if err != nil {
			continue
		}

		if ttl > 0 {
			totalTTL += ttl
			validKeys++
		}
	}

	if validKeys > 0 {
		stats.AverageTTL = totalTTL / time.Duration(validKeys)
	}

	return stats, nil
}

// APQStats represents APQ cache statistics.
type APQStats struct {
	TotalQueries int
	CacheSize    int
	AverageTTL   time.Duration
}

