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

type cachedResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (ca *Cache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		key := r.URL.String()
		if cached, found := ca.c.Get(key); found {
			resp := cached.(cachedResponse)
			for k, vv := range resp.Header {
				for _, v := range vv {
					w.Header().Add(k, v)
				}
			}
			w.Header().Set("Cache-Control", "max-age=300, HIT, from-cache=true")
			w.WriteHeader(resp.StatusCode)
			w.Write(resp.Body)
			return
		}

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		body := rec.Body.Bytes()
		resp := cachedResponse{
			StatusCode: rec.Code,
			Header:     rec.Header().Clone(),
			Body:       body,
		}
		if rec.Code == http.StatusOK {
			ca.c.Set(key, resp, cache.DefaultExpiration)
		}
		for k, vv := range rec.Header() {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		w.Header().Set("Cache-Control", "max-age=300, MISS, from-cache=false")
		w.WriteHeader(rec.Code)
		w.Write(body)
	})
}
