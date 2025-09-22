package balancer

import (
	"net/http"
	"sort"

	"github.com/dev-hak/mini-proxy/pkg/models"
)

type LeastConnections struct {
	*baseBalancer
}

func NewLeastConnections(backends []*models.Backend) *LeastConnections {
	return &LeastConnections{baseBalancer: &baseBalancer{backendsList: backends}}
}

func (lc *LeastConnections) NextBackend(req *http.Request) *models.Backend {
	backends := lc.backends()
	healthy := make([]*models.Backend, 0)
	for _, be := range backends {

		if be.Healthy {
			healthy = append(healthy, be)
		}
	}
	if len(healthy) == 0 {
		return nil
	}
	sort.Slice(healthy, func(i, j int) bool {
		return healthy[i].CurrentConns < healthy[j].CurrentConns
	})
	return healthy[0]
}

func (lc *LeastConnections) UpdateHealth(backend *models.Backend, healthy bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	for _, b := range lc.baseBalancer.backendsList {
		if b.URL == backend.URL {
			b.Healthy = healthy
			break
		}
	}
}
