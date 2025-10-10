package main

import (
	"log"
	"net/http"
	"time"

	"github.com/htooanttko/mini-proxy/internal/balancer"
	"github.com/htooanttko/mini-proxy/internal/config"
	"github.com/htooanttko/mini-proxy/internal/proxy/caching"
	"github.com/htooanttko/mini-proxy/internal/proxy/compression"
	"github.com/htooanttko/mini-proxy/internal/proxy/forward"
	"github.com/htooanttko/mini-proxy/internal/proxy/reverse"
	"github.com/htooanttko/mini-proxy/internal/proxy/security"
	"github.com/htooanttko/mini-proxy/internal/proxy/ssl"
	"github.com/htooanttko/mini-proxy/pkg/models"
)

func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	backends := make([]*models.Backend, len(cfg.Backends))
	for i := range cfg.Backends {
		backends[i] = &cfg.Backends[i]
	}
	bal := balancer.NewBalancer(backends, cfg.BalancerType)

	interval := time.Duration(cfg.HealthCheck.Interval)
	if interval <= 0 {
		log.Fatalf("Health check interval must be positive, got: %v", interval)
	}
	hc := balancer.NewHealthChecker(bal, interval, cfg.HealthCheck.Path)

	var handler http.Handler
	if cfg.ProxyType == "reverse" {
		rp := reverse.NewReverseProxy(bal, hc)
		handler = http.HandlerFunc(rp.ServeHTTP)
	} else {
		handler = http.HandlerFunc(forward.NewForwardProxy().ServeHTTP)
	}

	// Add middleware directly
	handler = security.NewRateLimiter(cfg.Security.RateLimit.RequestsPerMin).Middleware(handler)
	handler = security.NewAuth(cfg.Security.BasicAuthUsers).Middleware(handler)
	handler = caching.NewCache(time.Duration(cfg.Cache.DefaultExpiration), time.Duration(cfg.Cache.CleanupInterval)).Middleware(handler)

	handler = compression.Middleware(handler)

	// Middleware already in listener, but for SSL
	sslTerm := ssl.NewSSLTerminator(cfg.TLSCertFile, cfg.TLSKeyFile, handler)

	log.Printf("Starting server on %s", cfg.ListenAddr)
	if err := sslTerm.ListenAndServe(cfg.ListenAddr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	hc.Stop()
}
