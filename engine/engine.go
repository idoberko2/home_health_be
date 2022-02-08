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
		return ErrNotInitialized
	}

	go e.scheduler.Start(ctx, e.errReporter)
	return nil
}

func (e *engine) Ping(passphrase string) error {
	return e.healthCheck.Ping(passphrase)
}

func (e *engine) CheckState() (general.State, error) {
	isHealthy, err := e.healthCheck.IsHealthy()
	if err != nil {
		return general.StateUnhealthy, err
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

	return e.notifier.NotifyStateChange(newState)
}

var ErrNotInitialized = errors.New("not initialized")
