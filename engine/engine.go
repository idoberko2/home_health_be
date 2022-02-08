package engine

import (
	"context"
	"errors"

	"github.com/idoberko2/home_health_be/general"
	"github.com/idoberko2/home_health_be/healthcheck"
	"github.com/idoberko2/home_health_be/notifier"
	"github.com/idoberko2/home_health_be/scheduler"
)

type Engine interface {
	Init() error
	Start(ctx context.Context) error
}

func NewEngine(cfg EngineConfig, notifier notifier.Notifier) Engine {
	return &engine{
		cfg:      cfg,
		notifier: notifier,
	}
}

type engine struct {
	cfg         EngineConfig
	healthCheck healthcheck.HealthCheck
	scheduler   scheduler.Scheduler
	notifier    notifier.Notifier
	state       general.State
	errReporter chan error
	ready       bool
}

func (e *engine) Init() error {
	e.healthCheck = healthcheck.NewHealthCheck(e.cfg.HealthCheckConfig)
	e.scheduler = scheduler.NewScheduler(e.cfg.SchedulerConfig, e)
	e.ready = true

	return nil
}

func (e *engine) Start(ctx context.Context) error {
	if !e.ready {
		return errNotInitialized
	}

	go e.scheduler.Start(ctx, e.errReporter)
	return nil
}

func (e *engine) CheckState() (general.State, error) {
	return general.StateHealthy, nil
}

func (e *engine) OnStateCheck(general.State) error {
	return nil
}

var errNotInitialized = errors.New("not initialized")
