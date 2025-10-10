package balancer

import (
	"hash/fnv"
	"net/http"

	"github.com/htooanttko/mini-proxy/pkg/models"
)

type IPHash struct {
	*baseBalancer
}

func NewIPHash(backends []*models.Backend) *IPHash {
	return &IPHash{baseBalancer: &baseBalancer{backendsList: backends}}
}

func (ih *IPHash) NextBackend(req *http.Request) *models.Backend {
	ip := req.RemoteAddr
	h := fnv.New32a()
	h.Write([]byte(ip))
	hash := int(h.Sum32()) % len(ih.backends())
	be := ih.backends()[hash]
	if !be.Healthy {
		// Fallback to first healthy
		for _, b := range ih.backends() {
			if b.Healthy {
				return b
			}
		}
	}
	return be
}

func (ih *IPHash) UpdateHealth(backend *models.Backend, healthy bool) {
	ih.mu.Lock()
	defer ih.mu.Unlock()
	for _, b := range ih.baseBalancer.backendsList {
		if b.URL == backend.URL {
			b.Healthy = healthy
			break
		}
	}
}
