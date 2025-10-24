package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	limiter := &rateLimiter{
		requests: make(map[string]map[string][]time.Time),
	}

	// Make 10 requests with the same client IP and route.
	for range 10 {
		if !limiter.Allow("127.0.0.1", "/test", 10) {
			t.Fatalf("expected true got false")
		}
	}

	// Make the 11th request, should be denied.
	if limiter.Allow("127.0.0.1", "/test", 10) {
		t.Fatalf("expected false got true")
	}

}

func TestRateLimiter_WithDifferentEndpoints(t *testing.T) {
	limiter := &rateLimiter{
		requests: make(map[string]map[string][]time.Time),
	}

	// Make 10 requests with the same client IP and different routes.
	for range 10 {
		if !limiter.Allow("127.0.0.1", "/test1", 10) {
			t.Fatalf("expected true got false")
		}
	}

	if !limiter.Allow("127.0.0.1", "/test2", 100) {
		t.Fatalf("Different endpoints should not be affected")
	}

}

func TestRateLimiter_WithDifferentIPs(t *testing.T) {
	limiter := &rateLimiter{
		requests: make(map[string]map[string][]time.Time),
	}

	// Make 10 requests with different client IPs and the same route.
	for range 10 {
		if !limiter.Allow("127.0.0.1", "/test", 10) {
			t.Fatalf("expected true got false")
		}
	}

	if !limiter.Allow("127.0.0.2", "/test", 10) {
		t.Fatalf("Different IPs should not be affected")
	}

	if limiter.Allow("127.0.0.1", "/test", 10) {
		t.Fatalf("Same IP should be affected")
	}
}

func TestRateLimiter_Allo_TimeWindow(t *testing.T) {
	limiter := &rateLimiter{
		requests: make(map[string]map[string][]time.Time),
	}

	originalTime := time.Now()

	requestsPerMinute := 5

	testTimes := []time.Time{
		originalTime.Add(-30 * time.Second), // 30 seconds ago
		originalTime.Add(-20 * time.Second), // 20 seconds ago
		originalTime.Add(-10 * time.Second), // 10 seconds ago

	}

	// Manually add old requests
	limiter.requests["/test"] = map[string][]time.Time{
		"127.0.0.1": testTimes,
	}

	for range requestsPerMinute - len(testTimes) {
		if !limiter.Allow("127.0.0.1", "/test", requestsPerMinute) {
			t.Fatalf("expected true got false")
		}
	}

	// It should not allow new request
	if limiter.Allow("127.0.0.1", "/test", requestsPerMinute) {
		t.Fatalf("expected false got true")
	}

}
