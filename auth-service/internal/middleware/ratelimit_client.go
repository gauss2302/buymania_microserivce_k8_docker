package middleware

import (
	"net/http"
	"sync"
	"time"
)

type ClientRateLimiter struct {
	clients map[string]*RateLimiter
	mutex   sync.RWMutex
	config  ClientRateLimiterConfig
}

type ClientRateLimiterConfig struct {
	Capacity        float64
	RefillRate      float64
	CleanupInterval time.Duration
}

func NewClientRateLimiter(config ClientRateLimiterConfig) *ClientRateLimiter {
	limiter := &ClientRateLimiter{
		clients: make(map[string]*RateLimiter),
		config:  config,
	}

	// Start cleanup routine
	go limiter.cleanup()

	return limiter
}

func (crl *ClientRateLimiter) cleanup() {
	ticker := time.NewTicker(crl.config.CleanupInterval)
	for range ticker.C {
		crl.mutex.Lock()
		for ip, limiter := range crl.clients {
			if time.Since(limiter.lastRefill) > crl.config.CleanupInterval {
				delete(crl.clients, ip)
			}
		}
		crl.mutex.Unlock()
	}
}

func (crl *ClientRateLimiter) getLimiter(clientIP string) *RateLimiter {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	limiter, exists := crl.clients[clientIP]
	if !exists {
		limiter = NewRateLimiter(crl.config.Capacity, crl.config.RefillRate)
		crl.clients[clientIP] = limiter
	}
	return limiter
}

func PerClientRateLimit(limiter *ClientRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				clientIP = forwardedFor
			}

			if !limiter.getLimiter(clientIP).allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
