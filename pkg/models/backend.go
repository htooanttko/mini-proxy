package models

import (
	"time"
)

type Backend struct {
	URL             string    `json:"url"`
	Weight          int       `json:"weight"`
	MaxConnections  int       `json:"max_connections"`
	CurrentConns    int       `json:"current_conns"`
	Healthy         bool      `json:"healthy"`
	LastHealthCheck time.Time `json:"last_health_check"`
}
