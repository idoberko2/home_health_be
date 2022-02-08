package healthcheck

import (
	"errors"
	"time"
)

type HealthCheck interface {
	Ping(passphrase string) error
	IsHealthy() (bool, error)
	GetHistory() ([]time.Time, error)
}

func NewHealthCheck(cfg HealthCheckConfig) HealthCheck {
	return &healthcheck{
		cfg:         cfg,
		pingHistory: make([]time.Time, cfg.HistoryLength),
	}
}

type healthcheck struct {
	cfg         HealthCheckConfig
	pingHistory []time.Time
}

func (e *healthcheck) Ping(passphrase string) error {
	if passphrase != e.cfg.Passphrase {
		return ErrIncorrectPassphrase
	}

	// always enqueue the last ping in the head and maintain the queue at desired length
	e.pingHistory = append([]time.Time{time.Now()}, e.pingHistory[:e.cfg.HistoryLength-1]...)

	return nil
}

func (e *healthcheck) IsHealthy() (bool, error) {
	lastPing := e.pingHistory[0]
	if lastPing.IsZero() {
		return false, ErrNeverPinged
	}

	if lastPing.Add(e.cfg.GracePeriod).Before(time.Now()) {
		return false, nil
	}

	return true, nil
}

func (e *healthcheck) GetHistory() ([]time.Time, error) {
	return e.pingHistory, nil
}

var ErrIncorrectPassphrase = errors.New("not allowed")
var ErrNeverPinged = errors.New("never pinged")
