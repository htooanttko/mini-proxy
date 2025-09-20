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
	// log.Println("LeastConnections: Selecting next backend")
	backends := lc.backends()
	// log.Printf("LeastConnections: Available backends: %v", backends)
	healthy := make([]*models.Backend, 0)
	for _, be := range backends {
		// log.Printf("LeastConnections: Checking backend Weight %d", be.Weight)
		// log.Printf("LeastConnections: Checking backend MaxConnections %d", be.MaxConnections)
		// log.Printf("LeastConnections: Checking backend CurrentConns %d", be.CurrentConns)
		// log.Printf("LeastConnections: Checking backend Healthy %t", be.Healthy)
		// log.Printf("LeastConnections: Checking backend LastHealthCheck %d", be.LastHealthCheck)

		if be.Healthy {
			// log.Printf("LeastConnections: Backend %s is healthy", be.URL)
			healthy = append(healthy, be)
		} else {
			// log.Printf("LeastConnections: Backend %s is unhealthy", be.URL)
		}
	}
	if len(healthy) == 0 {
		// log.Println("LeastConnections: No healthy backends available")
		return nil
	}
	sort.Slice(healthy, func(i, j int) bool {
		return healthy[i].CurrentConns < healthy[j].CurrentConns
	})
	// log.Printf("LeastConnections: Selected backend %s with %d connections", healthy[0].URL, healthy[0].CurrentConns)
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
