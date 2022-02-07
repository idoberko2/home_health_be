package main

import "time"

type HealthCheckConfig struct {
	HistoryLength int
	Passphrase    string
	GracePeriod   time.Duration
}
