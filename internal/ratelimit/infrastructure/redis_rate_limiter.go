// Package ratelimit provides rate limiting infrastructure implementations.
package ratelimit

import (
	"context"
	"fmt"
	"time"

	ratelimitdomain "github.com/Rusty-Alucard/sp_api/internal/ratelimit/domain"
	"github.com/redis/go-redis/v9"
)

const (
	// SecondsPerMinute represents seconds in a minute
	SecondsPerMinute = 60
	// SecondsPerHour represents seconds in an hour
	SecondsPerHour   = 3600
	// SecondsPerDay represents seconds in a day
	SecondsPerDay    = 86400
)

// RedisRateLimiter implements rate limiting using Redis.
type RedisRateLimiter struct {
	client redis.UniversalClient
}

// NewRedisRateLimiter creates a new RedisRateLimiter.
func NewRedisRateLimiter(client redis.UniversalClient) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
	}
}

// Allow checks if a request is allowed under the rate limit.
func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit ratelimitdomain.RateLimit) (bool, error) {
	windows := limit.GetWindows()

	for _, window := range windows {
		allowed, err := r.checkWindow(ctx, key, window)
		if err != nil {
			return false, err
		}
		if !allowed {
			return false, nil
		}
	}

	return true, nil
}

// GetRemaining returns the number of remaining requests for a key.
func (r *RedisRateLimiter) GetRemaining(ctx context.Context, key string, limit ratelimitdomain.RateLimit) (int, error) {
	windows := limit.GetWindows()
	minRemaining := limit.RequestsPerMinute // Start with the most restrictive limit

	for _, window := range windows {
		remaining, err := r.getWindowRemaining(ctx, key, window)
		if err != nil {
			return 0, err
		}
		if remaining < minRemaining {
			minRemaining = remaining
		}
	}

	return minRemaining, nil
}

// GetResetTime returns the reset time for the most restrictive window.
func (r *RedisRateLimiter) GetResetTime(
	ctx context.Context,
	key string,
	limit ratelimitdomain.RateLimit,
) (time.Time, error) {
	windows := limit.GetWindows()
	earliestReset := time.Now().Add(ratelimitdomain.HoursPerDay * time.Hour) // Start with the longest window

	for _, window := range windows {
		resetTime, err := r.getWindowResetTime(ctx, key, window)
		if err != nil {
			return time.Time{}, err
		}
		if resetTime.Before(earliestReset) {
			earliestReset = resetTime
		}
	}

	return earliestReset, nil
}

// checkWindow checks if a request is allowed within a specific time window.
func (r *RedisRateLimiter) checkWindow(ctx context.Context, key string, window ratelimitdomain.Window) (bool, error) {
	windowKey := r.getWindowKey(key, window.Duration)

	// Use Redis pipeline for atomic operations
	pipe := r.client.Pipeline()

	// Increment counter
	incrCmd := pipe.Incr(ctx, windowKey)

	// Set expiration if this is the first request in the window
	pipe.Expire(ctx, windowKey, window.Duration)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := incrCmd.Val()
	return count <= int64(window.Limit), nil
}

// getWindowRemaining returns the remaining requests for a specific window.
func (r *RedisRateLimiter) getWindowRemaining(
	ctx context.Context,
	key string,
	window ratelimitdomain.Window,
) (int, error) {
	windowKey := r.getWindowKey(key, window.Duration)

	count, err := r.client.Get(ctx, windowKey).Int()
	if err == redis.Nil {
		// Key doesn't exist, so all requests are remaining
		return window.Limit, nil
	}
	if err != nil {
		return 0, err
	}

	remaining := window.Limit - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// getWindowResetTime returns the reset time for a specific window.
func (r *RedisRateLimiter) getWindowResetTime(
	ctx context.Context,
	key string,
	window ratelimitdomain.Window,
) (time.Time, error) {
	windowKey := r.getWindowKey(key, window.Duration)

	ttl, err := r.client.TTL(ctx, windowKey).Result()
	if err != nil {
		return time.Time{}, err
	}

	if ttl == -1 {
		// Key exists but has no expiration, set it
		err = r.client.Expire(ctx, windowKey, window.Duration).Err()
		if err != nil {
			return time.Time{}, err
		}
		ttl = window.Duration
	}

	if ttl == -2 {
		// Key doesn't exist
		return time.Now().Add(window.Duration), nil
	}

	return time.Now().Add(ttl), nil
}

// getWindowKey generates a Redis key for a specific time window.
func (r *RedisRateLimiter) getWindowKey(key string, duration time.Duration) string {
	now := time.Now()

	switch duration {
	case time.Nanosecond, time.Microsecond, time.Millisecond, time.Second:
		return fmt.Sprintf("rate_limit:%s:custom:%d", key, now.Unix()/int64(duration.Seconds()))
	case time.Minute:
		return fmt.Sprintf("rate_limit:%s:minute:%d", key, now.Unix()/SecondsPerMinute)
	case time.Hour:
		return fmt.Sprintf("rate_limit:%s:hour:%d", key, now.Unix()/SecondsPerHour)
	case ratelimitdomain.HoursPerDay * time.Hour:
		return fmt.Sprintf("rate_limit:%s:day:%d", key, now.Unix()/SecondsPerDay)
	default:
		return fmt.Sprintf("rate_limit:%s:custom:%d", key, now.Unix()/int64(duration.Seconds()))
	}
}
