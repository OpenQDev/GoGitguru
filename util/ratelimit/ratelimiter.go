package ratelimit

import (
	"sync"
	"time"
)

// RequestInfo stores the timestamps of recent requests for a client.
type RequestInfo struct {
	Timestamps []time.Time
	mu         sync.Mutex
}

// RateLimiter controls access to the rate-limited resources.
type RateLimiter struct {
	requests map[string]*RequestInfo
	limit    int
	window   time.Duration
	mu       sync.Mutex
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]*RequestInfo),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if a request from the given IP is allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if _, exists := rl.requests[ip]; !exists {
		rl.requests[ip] = &RequestInfo{}
	}

	reqInfo := rl.requests[ip]
	reqInfo.mu.Lock()
	defer reqInfo.mu.Unlock()

	// Remove old timestamps outside of the current window.
	validAfter := now.Add(-rl.window)
	newTimestamps := make([]time.Time, 0)
	for _, timestamp := range reqInfo.Timestamps {
		if timestamp.After(validAfter) {
			newTimestamps = append(newTimestamps, timestamp)
		}
	}
	reqInfo.Timestamps = newTimestamps

	// Check if adding a new request would exceed the limit.
	if len(reqInfo.Timestamps) >= rl.limit {
		return false
	}

	// Add the new request timestamp.
	reqInfo.Timestamps = append(reqInfo.Timestamps, now)
	return true
}
