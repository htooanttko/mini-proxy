package balancer

import (
	"net/http"
	"sync"

	"github.com/dev-hak/mini-proxy/pkg/models"
)

type Balancer interface {
	NextBackend(*http.Request) *models.Backend
	UpdateHealth(*models.Backend, bool)
}

type baseBalancer struct {
	mu           sync.RWMutex
	backendsList []*models.Backend
}

func NewBalancer(backends []*models.Backend, balancerType string) Balancer {
	switch balancerType {
	case "round_robin":
		return NewRoundRobin(backends)
	case "least_conn":
		return NewLeastConnections(backends)
	case "ip_hash":
		return NewIPHash(backends)
	case "weighted":
		return NewWeighted(backends)
	default:
		return NewRoundRobin(backends)
	}
}

func (b *baseBalancer) backends() []*models.Backend {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.backendsList
}

func (b *baseBalancer) NextBackend(req *http.Request) *models.Backend {
	// Placeholder implementation
	return nil
}

func (b *baseBalancer) UpdateHealth(backend *models.Backend, healthy bool) {
	// Placeholder implementation
}
