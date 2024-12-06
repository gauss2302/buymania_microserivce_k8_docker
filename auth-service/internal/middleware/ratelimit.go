package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens        float64
	capacity      float64
	refillRate    float64
	lastRefill    time.Time
	mutex         sync.Mutex
	perClientRate map[string]float64 // Track rates per IP
}

func NewRateLimiter(capacity float64, refillRate float64) *RateLimiter {
	return &RateLimiter{
		tokens:        capacity,
		capacity:      capacity,
		refillRate:    refillRate,
		lastRefill:    time.Now(),
		perClientRate: make(map[string]float64),
	}
}

func (rl *RateLimiter) refill() {
	now := time.Now()
	duration := now.Sub(rl.lastRefill).Seconds()
	rl.tokens = min(rl.capacity, rl.tokens+(duration*rl.refillRate))
	rl.lastRefill = now
}

func (rl *RateLimiter) allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	rl.refill()
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

func RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
