package balancer

import (
	"net/http"
	"sort"

	"github.com/htooanttko/mini-proxy/pkg/models"
)

type Weighted struct {
	*baseBalancer
}

func NewWeighted(backends []*models.Backend) *Weighted {
	return &Weighted{baseBalancer: &baseBalancer{backendsList: backends}}
}

func (w *Weighted) NextBackend(req *http.Request) *models.Backend {
	healthy := make([]*models.Backend, 0)
	totalWeight := 0
	for _, be := range w.backends() {
		if be.Healthy {
			healthy = append(healthy, be)
			totalWeight += be.Weight
		}
	}
	if len(healthy) == 0 {
		return nil
	}
	// Simple weighted selection (can be improved with reservoir sampling)
	r := sort.Search(len(healthy), func(i int) bool {
		sum := 0
		for j := 0; j <= i; j++ {
			sum += healthy[j].Weight
		}
		return sum >= totalWeight/2 // Pseudo-random
	})
	return healthy[r%len(healthy)]
}

func (w *Weighted) UpdateHealth(backend *models.Backend, healthy bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, b := range w.baseBalancer.backendsList {
		if b.URL == backend.URL {
			b.Healthy = healthy
			break
		}
	}
}
