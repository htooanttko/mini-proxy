package security

import (
	"net/http"
	"sync"

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
