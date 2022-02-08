package engine

import (
	"context"
	"testing"
	"time"

	"github.com/idoberko2/home_health_be/healthcheck"
	"github.com/idoberko2/home_health_be/scheduler"
	"github.com/stretchr/testify/assert"
)

const somePassphrase = "somePassphrase"

func testEngineConfig() EngineConfig {
	return EngineConfig{
		HealthCheckConfig: healthcheck.HealthCheckConfig{
			HistoryLength: 10,
			Passphrase:    somePassphrase,
			GracePeriod:   time.Second,
		},
		SchedulerConfig: scheduler.SchedulerConfig{
			SampleRate: 100 * time.Millisecond,
		},
	}
}

func TestStartNoInit(t *testing.T) {
	engine := NewEngine(testEngineConfig(), nil)
	err := engine.Start(context.Background())

	assert.Equal(t, ErrNotInitialized, err)
}
