package reverse

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dev-hak/mini-proxy/internal/balancer"
)

type ReverseProxy struct {
	balancer balancer.Balancer
	health   *balancer.HealthChecker
}

func NewReverseProxy(bal balancer.Balancer, hc *balancer.HealthChecker) *ReverseProxy {
	return &ReverseProxy{balancer: bal, health: hc}
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	be := rp.balancer.NextBackend(req)
	if be == nil {
		http.Error(w, "No healthy backend", http.StatusServiceUnavailable)
		return
	}

	target, _ := url.Parse(be.URL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify request
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.Host = target.Host

	// Increment connections
	be.CurrentConns++
	defer func() { be.CurrentConns-- }()

	proxy.ServeHTTP(w, req)
}
