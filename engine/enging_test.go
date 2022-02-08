package engine

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/idoberko2/home_health_be/general"
	"github.com/idoberko2/home_health_be/healthcheck"
	"github.com/idoberko2/home_health_be/scheduler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const somePassphrase = "somePassphrase"

func testEngineConfig() EngineConfig {
	return EngineConfig{
		HealthCheckConfig: healthcheck.HealthCheckConfig{
			HistoryLength: 10,
			Passphrase:    somePassphrase,
			GracePeriod:   50 * time.Millisecond,
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

func TestHealthNoPing(t *testing.T) {
	notif := &notifierMock{}
	notif.On("NotifyStateChange", mock.Anything).Return(nil)
	engine := NewEngine(testEngineConfig(), notif)
	err := engine.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(110*time.Millisecond))
	defer cancel()
	err = engine.Start(ctx)
	assert.NoError(t, err)

	<-ctx.Done()
	notif.AssertNotCalled(t, "NotifyStateChange", mock.Anything)
}

func TestHealthyAlways(t *testing.T) {
	notif := &notifierMock{}
	notif.On("NotifyStateChange", mock.Anything).Return(nil)
	engine := NewEngine(testEngineConfig(), notif)
	err := engine.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(200*time.Millisecond))
	defer cancel()
	err = engine.Start(ctx)
	assert.NoError(t, err)

	err = engine.Ping(somePassphrase)
	assert.NoError(t, err)

	<-ctx.Done()
	notif.AssertCalled(t, "NotifyStateChange", general.StateHealthy)
}

type notifierMock struct {
	mock.Mock
}

func (n *notifierMock) NotifyStateChange(state general.State) error {
	log.Printf("notify state change %d\n", state)
	args := n.Called(state)

	return args.Error(0)
}
