package ratelimit

import (
	"context"
	"errors"
	"time"
)

// RateLimit represents rate limiting configuration.
type RateLimit struct {
	RequestsPerMinute int
	RequestsPerHour   int
	RequestsPerDay    int
}

// Window represents a time window for rate limiting.
type Window struct {
	Duration time.Duration
	Limit    int
}

// RateLimitError represents a rate limit exceeded error.
type RateLimitError struct {
	Message    string
	RetryAfter time.Duration
	Limit      int
	Remaining  int
	ResetTime  time.Time
}

// Error implements the error interface.
func (e *RateLimitError) Error() string {
	return e.Message
}

// GetWindows returns all time windows for rate limiting.
func (r RateLimit) GetWindows() []Window {
	return []Window{
		{Duration: time.Minute, Limit: r.RequestsPerMinute},
		{Duration: time.Hour, Limit: r.RequestsPerHour},
		{Duration: 24 * time.Hour, Limit: r.RequestsPerDay},
	}
}

// IsRateLimitExceeded checks if an error is a rate limit exceeded error.
func IsRateLimitExceeded(err error) bool {
	var rateLimitErr *RateLimitError
	return errors.As(err, &rateLimitErr)
}

// NewRateLimitExceededError creates a new rate limit exceeded error.
func NewRateLimitExceededError(retryAfter time.Duration, limit, remaining int, resetTime time.Time) *RateLimitError {
	return &RateLimitError{
		Message:    "rate limit exceeded",
		RetryAfter: retryAfter,
		Limit:      limit,
		Remaining:  remaining,
		ResetTime:  resetTime,
	}
}

// RateLimiter defines the interface for rate limiting operations.
type RateLimiter interface {
	// Allow checks if a request is allowed under the rate limit.
	Allow(ctx context.Context, key string, limit RateLimit) (bool, error)

	// GetRemaining returns the number of remaining requests for a key.
	GetRemaining(ctx context.Context, key string, limit RateLimit) (int, error)

	// GetResetTime returns the reset time for the most restrictive window.
	GetResetTime(ctx context.Context, key string, limit RateLimit) (time.Time, error)
}
