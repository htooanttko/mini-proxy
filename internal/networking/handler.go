package networking

import (
	"net/http"

	reverse_proxy "github.com/dev-hak/mini-proxy/internal/proxy/reverse"
)

func ProxyHandler(rp *reverse_proxy.ReverseProxy) http.HandlerFunc {
	return rp.ServeHTTP
}
