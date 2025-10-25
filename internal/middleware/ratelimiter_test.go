package middleware

import (
	"testing"
)

func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter()
	defer limiter.Close()

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
	limiter := NewRateLimiter()
	defer limiter.Close()

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
	limiter := NewRateLimiter()
	defer limiter.Close()

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
	limiter := NewRateLimiter()
	defer limiter.Close()

	requestsPerMinute := 5

	// Make requests up to the limit
	for range requestsPerMinute {
		if !limiter.Allow("127.0.0.1", "/test", requestsPerMinute) {
			t.Fatalf("expected true got false")
		}
	}

	// Next request should be denied
	if limiter.Allow("127.0.0.1", "/test", requestsPerMinute) {
		t.Fatalf("expected false got true")
	}
}
