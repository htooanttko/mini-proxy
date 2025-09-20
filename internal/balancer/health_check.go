package balancer

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type HealthChecker struct {
	balancers []Balancer
	interval  time.Duration
	ticker    *time.Ticker
	done      chan struct{}
	mu        sync.Mutex
}

func (hc *HealthChecker) Interval() time.Duration {
	return hc.interval
}

func NewHealthChecker(bal Balancer, interval time.Duration) *HealthChecker {
	hc := &HealthChecker{
		balancers: []Balancer{bal},
		interval:  interval,
		ticker:    time.NewTicker(interval),
		done:      make(chan struct{}),
	}
	hc.checkAll()

	go hc.runCheckLoop()
	return hc
}

func (hc *HealthChecker) runCheckLoop() {
	for {
		select {
		case <-hc.ticker.C:
			hc.checkAll()
		case <-hc.done:
			hc.ticker.Stop()
			return
		}
	}
}

func (hc *HealthChecker) checkAll() {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	for _, bal := range hc.balancers {
		switch b := bal.(type) {
		case *RoundRobin:
			hc.checkBackend(b.baseBalancer)
		case *LeastConnections:
			hc.checkBackend(b.baseBalancer)
		case *IPHash:
			hc.checkBackend(b.baseBalancer)
		case *Weighted:
			hc.checkBackend(b.baseBalancer)
		default:
			log.Printf("Balancer type %T does not support health checks", bal)
		}
	}
}

func (hc *HealthChecker) checkBackend(base *baseBalancer) {
	for _, be := range base.backendsList {
		resp, err := http.Get(be.URL + "/health")
		if err != nil {
			be.Healthy = false
		} else {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				be.Healthy = true
			} else {
				be.Healthy = false
			}
			resp.Body.Close()
		}
	}
}

func (hc *HealthChecker) Stop() {
	close(hc.done)
}
