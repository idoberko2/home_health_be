package engine

import (
	"context"

	"github.com/idoberko2/home_health_be/general"
	"github.com/idoberko2/home_health_be/healthcheck"
	"github.com/idoberko2/home_health_be/scheduler"
)

type Engine interface {
	Init() error
	Start(ctx context.Context) error
}

func NewEngine(cfg EngineConfig) Engine {
	return &engine{
		cfg: cfg,
	}
}

type engine struct {
	cfg         EngineConfig
	healthCheck healthcheck.HealthCheck
	scheduler   scheduler.Scheduler
	state       general.State
	errReporter chan error
}

func (e *engine) Init() error {
	e.healthCheck = healthcheck.NewHealthCheck(e.cfg.HealthCheckConfig)
	e.scheduler = scheduler.NewScheduler(e.cfg.SchedulerConfig, e)

	return nil
}

func (e *engine) Start(ctx context.Context) error {
	go e.scheduler.Start(ctx, e.errReporter)
	return nil
}

func (e *engine) CheckState() (general.State, error) {
	return general.StateHealthy, nil
}

func (e *engine) OnStateCheck(general.State) error {
	return nil
}
