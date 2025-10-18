package balancer

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/htooanttko/mini-proxy/pkg/models"
)

type RoundRobin struct {
	*baseBalancer
	index int64
}

func NewRoundRobin(backends []*models.Backend) *RoundRobin {
	return &RoundRobin{
		baseBalancer: &baseBalancer{backendsList: backends},
	}
}

func (rr *RoundRobin) NextBackend(req *http.Request) *models.Backend {
	backends := rr.backends()
	if len(backends) == 0 {
		return nil
	}
	// Filter healthy
	healthy := make([]*models.Backend, 0)
	for _, be := range backends {
		if be.Healthy {
			healthy = append(healthy, be)
		}
	}
	if len(healthy) == 0 {
		return nil
	}
	idx := atomic.AddInt64(&rr.index, 1) % int64(len(healthy))
	return healthy[idx]
}

func (rr *RoundRobin) UpdateHealth(backend *models.Backend, healthy bool) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for _, b := range rr.baseBalancer.backendsList {
		if b.URL == backend.URL {
			b.Healthy = healthy
			break
		}
	}
}
