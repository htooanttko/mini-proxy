package balancer

import (
	"log"
	"net/http"
	"sync"

	"github.com/htooanttko/mini-proxy/pkg/models"
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
	b.mu.Lock()
	defer b.mu.Unlock()

	log.Printf("Selecting backend from list: %v", b.backendsList)
	// Example: Round-robin logic
	if len(b.backendsList) == 0 {
		return nil
	}
	backend := b.backendsList[0]
	b.backendsList = append(b.backendsList[1:], backend)
	return backend
}

func (b *baseBalancer) UpdateHealth(backend *models.Backend, healthy bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, existingBackend := range b.backendsList {
		if existingBackend.URL == backend.URL {
			existingBackend.Healthy = healthy
			if !healthy {
				// Remove unhealthy backend
				b.backendsList = append(b.backendsList[:i], b.backendsList[i+1:]...)
			}
			break
		}
	}
}
