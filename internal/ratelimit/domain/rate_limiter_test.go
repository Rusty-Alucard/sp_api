package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimit_GetWindows(t *testing.T) {
	rateLimit := RateLimit{
		RequestsPerMinute: 10,
		RequestsPerHour:   100,
		RequestsPerDay:    1000,
	}

	windows := rateLimit.GetWindows()

	assert.Len(t, windows, 3)

	// 分単位のウィンドウ
	assert.Equal(t, time.Minute, windows[0].Duration)
	assert.Equal(t, 10, windows[0].Limit)

	// 時間単位のウィンドウ
	assert.Equal(t, time.Hour, windows[1].Duration)
	assert.Equal(t, 100, windows[1].Limit)

	// 日単位のウィンドウ
	assert.Equal(t, HoursPerDay*time.Hour, windows[2].Duration)
	assert.Equal(t, 1000, windows[2].Limit)
}

func TestRateLimitError_Error(t *testing.T) {
	err := &Error{
		Message: "rate limit exceeded",
	}

	assert.Equal(t, "rate limit exceeded", err.Error())
}

func TestIsRateLimitExceeded(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "RateLimitError",
			err:      &Error{Message: "rate limit exceeded"},
			expected: true,
		},
		{
			name:     "Other error",
			err:      assert.AnError,
			expected: false,
		},
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRateLimitExceeded(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewRateLimitExceededError(t *testing.T) {
	retryAfter := 30 * time.Second
	limit := 100
	remaining := 0
	resetTime := time.Now().Add(time.Minute)

	err := NewRateLimitExceededError(retryAfter, limit, remaining, resetTime)

	assert.Equal(t, "rate limit exceeded", err.Message)
	assert.Equal(t, retryAfter, err.RetryAfter)
	assert.Equal(t, limit, err.Limit)
	assert.Equal(t, remaining, err.Remaining)
	assert.Equal(t, resetTime, err.ResetTime)
}
