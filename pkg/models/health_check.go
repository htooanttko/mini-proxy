package models

import "time"

type HealthCheck interface {
	Run()
	Stop()
	Interval() time.Duration
}
