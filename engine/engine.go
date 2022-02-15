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
	Ping(passphrase string) error
}

func New(cfg EngineConfig, notifier notifier.Notifier) Engine {
	return &engine{
		cfg:      cfg,
		notifier: notifier,
		state:    general.StateUndefined,
	}
}

type engine struct {
	cfg         EngineConfig
	healthCheck healthcheck.HealthCheck
	scheduler   scheduler.Scheduler
	notifier    notifier.Notifier
	state       general.State
	ready       bool
}

func (e *engine) Init() error {
	e.healthCheck = healthcheck.New(e.cfg.HealthCheckConfig)
	e.scheduler = scheduler.New(e.cfg.SchedulerConfig, e)
	e.ready = true

	return nil
}

func (e *engine) Start(ctx context.Context) error {
	if !e.ready {
		return ErrNotInitialized
	}

	e.scheduler.Start(ctx, e.cfg.errReporter)
	return nil
}

func (e *engine) Ping(passphrase string) error {
	return e.healthCheck.Ping(passphrase)
}

func (e *engine) CheckState() (general.State, error) {
	isHealthy, err := e.healthCheck.IsHealthy()
	if err != nil {
		return general.StateUndefined, err
	}

	if !isHealthy {
		return general.StateUnhealthy, nil
	}

	return general.StateHealthy, nil
}

func (e *engine) OnStateCheck(newState general.State) error {
	if newState == e.state {
		return nil
	}

	e.state = newState

	return e.notifier.NotifyStateChange(newState)
}

var ErrNotInitialized = errors.New("not initialized")
