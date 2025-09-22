package security

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      float64 // requests per second
}

func NewRateLimiter(requestsPerMin int) *RateLimiter {
	rps := float64(requestsPerMin) / 60
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		rl.mu.Lock()
		limiter, exists := rl.limiters[clientIP]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rl.rps), int(rl.rps*2)) // burst
			rl.limiters[clientIP] = limiter
		}
		rl.mu.Unlock()

		if !limiter.Allow() {
			http.Error(w, "Rate limited", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type clientRequestTracker struct {
	requests      map[string]int
	lastResetTime time.Time
	mutex         sync.Mutex
}

func NewClientRequestTracker() *clientRequestTracker {
	return &clientRequestTracker{
		requests:      make(map[string]int),
		lastResetTime: time.Now(),
	}
}

func (crt *clientRequestTracker) Middleware(limit int, resetInterval time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			crt.mutex.Lock()
			defer crt.mutex.Unlock()

			// Reset request counts periodically
			if time.Since(crt.lastResetTime) > resetInterval {
				crt.requests = make(map[string]int)
				crt.lastResetTime = time.Now()
			}

			// Check rate limit
			if crt.requests[clientIP] >= limit {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Increment request count
			crt.requests[clientIP]++
			next.ServeHTTP(w, r)
		})
	}
}
