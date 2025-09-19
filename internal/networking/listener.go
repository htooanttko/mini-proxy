package networking

import (
	"net/http"

	"github.com/dev-hak/mini-proxy/internal/config"
	"github.com/dev-hak/mini-proxy/internal/proxy/caching"
	"github.com/dev-hak/mini-proxy/internal/proxy/compression"
	forward_proxy "github.com/dev-hak/mini-proxy/internal/proxy/forward"
	"github.com/dev-hak/mini-proxy/internal/proxy/security"
)

type Listener struct {
	mux *http.ServeMux
}

func NewListener(cfg *config.Config) *Listener {
	mux := http.NewServeMux()
	// Add routes, e.g., mux.Handle("/", proxyHandler)

	// Middleware stack
	var handler http.Handler = mux
	handler = security.NewRateLimiter(cfg.Security.RateLimit.RequestsPerMin).Middleware(handler)
	handler = security.NewAuth(cfg.Security.BasicAuthUsers).Middleware(handler)
	handler = caching.NewCache(cfg.Cache.DefaultExpiration, cfg.Cache.CleanupInterval).Middleware(handler)
	handler = compression.NewCompressor().Middleware(handler)

	l := &Listener{mux: mux}
	l.setupHandlers(cfg, handler)
	return l
}

func (l *Listener) setupHandlers(cfg *config.Config, base http.Handler) {
	if cfg.ProxyType == "forward" {
		l.mux.Handle("/", http.HandlerFunc(forward_proxy.NewForwardProxy().ServeHTTP))
	} else {
		// Reverse
		l.mux.Handle("/", base) // Will be set in main
	}
}

func (l *Listener) Serve(addr string) error {
	return http.ListenAndServe(addr, l.mux)
}

func (l *Listener) ListenAndServe(handler http.Handler) error {
	srv := &http.Server{
		Addr:    ":8080", // Replace with dynamic config if needed
		Handler: handler,
	}
	return srv.ListenAndServe()
}
