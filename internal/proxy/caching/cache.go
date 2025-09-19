package caching

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	c *cache.Cache
}

func NewCache(exp time.Duration, cleanup time.Duration) *Cache {
	return &Cache{c: cache.New(exp, cleanup)}
}

func (ca *Cache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.String()
		if _, found := ca.c.Get(key); found {
			// Assume response is cached
			w.WriteHeader(200)
			w.Write([]byte("Cached response")) // Placeholder
			return
		}

		// Capture response
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		// Cache if cacheable
		if rec.Code == http.StatusOK && r.Method == "GET" {
			ca.c.Set(key, "Cached body", cache.DefaultExpiration)
		}

		// Write to original
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	})
}
