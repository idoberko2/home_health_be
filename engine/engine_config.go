package engine

import (
	"github.com/idoberko2/home_health_be/healthcheck"
	"github.com/idoberko2/home_health_be/scheduler"
)

type EngineConfig struct {
	healthcheck.HealthCheckConfig
	scheduler.SchedulerConfig
	errReporter chan error
}
