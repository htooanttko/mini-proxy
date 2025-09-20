package models

import (
	"encoding/json"
	"errors"
	"time"
)

type Backend struct {
	URL             string
	Weight          int
	MaxConnections  int
	CurrentConns    int
	Healthy         bool
	LastHealthCheck time.Time
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return errors.New("invalid duration format")
	}
	*d = Duration(dur)
	return nil
}

type HealthCheckConfig struct {
	Interval Duration `json:"Interval"`
	Timeout  Duration `json:"Timeout"`
	Path     string   `json:"Path"`
}
