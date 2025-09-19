package balancer

import (
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
		// Refactor type assertion to ensure compatibility with `Balancer` interface.
		if base, ok := bal.(*baseBalancer); ok {
			for _, be := range base.backendsList {
				resp, err := http.Get(be.URL + "/health")
				if err != nil {
					be.Healthy = false
				} else {
					resp.Body.Close()
					be.Healthy = true
				}
			}
		}
	}
}

func (hc *HealthChecker) Stop() {
	close(hc.done)
}
