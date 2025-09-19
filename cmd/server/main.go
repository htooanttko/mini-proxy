package main

import (
	"log"
	"net/http"

	"github.com/dev-hak/mini-proxy/internal/balancer"
	"github.com/dev-hak/mini-proxy/internal/config"
	"github.com/dev-hak/mini-proxy/internal/networking"

	forward_proxy "github.com/dev-hak/mini-proxy/internal/proxy/forward"
	reverse_proxy "github.com/dev-hak/mini-proxy/internal/proxy/reverse"
	"github.com/dev-hak/mini-proxy/internal/proxy/ssl"
	"github.com/dev-hak/mini-proxy/pkg/models"
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
	hc := balancer.NewHealthChecker(bal, cfg.HealthCheck.Interval)

	var handler http.Handler
	if cfg.ProxyType == "reverse" {
		rp := reverse_proxy.NewReverseProxy(bal, hc)
		handler = networking.ProxyHandler(rp)
	} else {
		// Forward proxy handler
		handler = http.HandlerFunc(forward_proxy.NewForwardProxy().ServeHTTP)
	}

	// Middleware already in listener, but for SSL
	sslTerm := ssl.NewSSLTerminator(cfg.TLSCertFile, cfg.TLSKeyFile, handler)

	listener := networking.NewListener(cfg)
	log.Printf("Starting server on %s", cfg.ListenAddr)
	if err := listener.ListenAndServe(http.HandlerFunc(sslTerm.ServeHTTP)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	hc.Stop()
}
