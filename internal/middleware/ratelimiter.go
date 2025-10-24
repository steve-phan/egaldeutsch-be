package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimitConfig struct {
	RequestsPerMinute int
}

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string]map[string][]time.Time // Endpoint -> Client IP -> Requests
}

func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		route := c.Request.URL.Path

		if !sharedRateLimiter.Allow(clientIP, route, config.RequestsPerMinute) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()

	}
}

// Global shared rate limiter instance
var sharedRateLimiter = &rateLimiter{
	requests: make(map[string]map[string][]time.Time),
}

// Start cleanup goroutine
func init() {
	go sharedRateLimiter.cleanupOldEntries()
}

func (r *rateLimiter) Allow(clientIP string, route string, limit int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Initialize nested maps if needed
	if r.requests[route] == nil {
		r.requests[route] = make(map[string][]time.Time)
	}

	now := time.Now()
	windowStart := now.Add(-time.Minute)

	// Get current requests and filter old ones
	requests := r.requests[route][clientIP]
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
	r.requests[route][clientIP] = recentRequests
	return true
}

func (r *rateLimiter) cleanupOldEntries() {
	for {
		time.Sleep(5 * time.Minute)
		r.mu.Lock()

		now := time.Now()
		windowStart := now.Add(-2 * time.Minute)

		for _, routeRequests := range r.requests {
			for clientId, requests := range routeRequests {
				var recentRequests []time.Time
				for _, request := range requests {
					if request.After(windowStart) {
						recentRequests = append(recentRequests, request)
					}
				}
				if len(recentRequests) == 0 {
					delete(routeRequests, clientId)
				} else {
					routeRequests[clientId] = recentRequests
				}
			}

		}

		r.mu.Unlock()
	}
}
