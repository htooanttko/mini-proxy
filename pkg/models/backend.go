package models

import "time"

type Backend struct {
	URL             string
	Weight          int
	MaxConnections  int
	CurrentConns    int
	Healthy         bool
	LastHealthCheck time.Time
}

type HealthCheckConfig struct {
	Interval time.Duration `json:"Interval"`
	Timeout  time.Duration `json:"Timeout"`
	Path     string        `json:"Path"`
}
