// Package ratelimit provides domain models and interfaces for rate limiting functionality.
package ratelimit

import (
	"context"
	"errors"
	"time"
)

const (
	// HoursPerDay represents the number of hours in a day.
	HoursPerDay = 24
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

// Error represents a rate limit exceeded error.
type Error struct {
	Message    string
	RetryAfter time.Duration
	Limit      int
	Remaining  int
	ResetTime  time.Time
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// GetWindows returns all time windows for rate limiting.
func (r RateLimit) GetWindows() []Window {
	return []Window{
		{Duration: time.Minute, Limit: r.RequestsPerMinute},
		{Duration: time.Hour, Limit: r.RequestsPerHour},
		{Duration: HoursPerDay * time.Hour, Limit: r.RequestsPerDay},
	}
}

// IsRateLimitExceeded checks if an error is a rate limit exceeded error.
func IsRateLimitExceeded(err error) bool {
	var rateLimitErr *Error
	return errors.As(err, &rateLimitErr)
}

// NewRateLimitExceededError creates a new rate limit exceeded error.
func NewRateLimitExceededError(retryAfter time.Duration, limit, remaining int, resetTime time.Time) *Error {
	return &Error{
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
