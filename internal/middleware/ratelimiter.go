package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter handles rate limiting logic.
// Following Go philosophy: explicit dependencies, no global state.
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]map[string][]time.Time // Endpoint -> Client IP -> Requests
	cleanup  chan struct{}                     // Channel to signal cleanup goroutine shutdown
}

// NewRateLimiter creates a new rate limiter instance.
// Explicit constructor following Go philosophy.
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]map[string][]time.Time),
		cleanup:  make(chan struct{}),
	}
	go rl.cleanupOldEntries()
	return rl
}

// Close stops the cleanup goroutine and releases resources.
func (rl *RateLimiter) Close() {
	close(rl.cleanup)
}

// Middleware returns a Gin middleware function that enforces rate limiting.
func (rl *RateLimiter) Middleware(requestsPerMinute int) gin.HandlerFunc {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 60 // sensible default
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		route := c.Request.URL.Path

		if !rl.Allow(clientIP, route, requestsPerMinute) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}

// Allow checks if a request should be allowed based on rate limiting rules.
func (rl *RateLimiter) Allow(clientIP string, route string, limit int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Initialize nested maps if needed
	if rl.requests[route] == nil {
		rl.requests[route] = make(map[string][]time.Time)
	}

	now := time.Now()
	windowStart := now.Add(-time.Minute)

	// Get current requests and filter old ones
	requests := rl.requests[route][clientIP]
	var recentRequests []time.Time

	for _, request := range requests {
		if request.After(windowStart) {
			recentRequests = append(recentRequests, request)
		}
	}

	// Check rate limit
	if len(recentRequests) >= limit {
		return false
	}

	// Add current request and update storage
	recentRequests = append(recentRequests, now)
	rl.requests[route][clientIP] = recentRequests
	return true
}

// cleanupOldEntries runs in a goroutine to periodically clean up old request records.
func (rl *RateLimiter) cleanupOldEntries() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.performCleanup()
		case <-rl.cleanup:
			return
		}
	}
}

// performCleanup removes old request records that are outside the time window.
func (rl *RateLimiter) performCleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-2 * time.Minute)

	for route, routeRequests := range rl.requests {
		for clientID, requests := range routeRequests {
			var recentRequests []time.Time
			for _, request := range requests {
				if request.After(windowStart) {
					recentRequests = append(recentRequests, request)
				}
			}
			if len(recentRequests) == 0 {
				delete(routeRequests, clientID)
			} else {
				routeRequests[clientID] = recentRequests
			}
		}

		// Remove empty route maps
		if len(routeRequests) == 0 {
			delete(rl.requests, route)
		}
	}
}

// Default global rate limiter instance for backward compatibility
var defaultRateLimiter = NewRateLimiter()

// RateLimit creates a simple rate limiting middleware with a global limiter.
// Convenience function for backward compatibility.
func RateLimit(limit int) gin.HandlerFunc {
	return defaultRateLimiter.Middleware(limit)
}
