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
			GracePeriod:   500 * time.Millisecond,
		},
		SchedulerConfig: scheduler.SchedulerConfig{
			SampleRate: 100 * time.Millisecond,
		},
	}
}

func TestStartNoInit(t *testing.T) {
	engine := New(testEngineConfig(), nil)
	err := engine.Start(context.Background())

	assert.Equal(t, ErrNotInitialized, err)
}

func TestHealthNoPing(t *testing.T) {
	notif := &notifierMock{}
	notif.On("NotifyStateChange", mock.Anything).Return(nil)
	engine := New(testEngineConfig(), notif)
	err := engine.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(110*time.Millisecond))
	defer cancel()
	go engine.Start(ctx)

	<-ctx.Done()
	notif.AssertNotCalled(t, "NotifyStateChange", mock.Anything)
}

func TestHealthyAlways(t *testing.T) {
	notif := &notifierMock{}
	engineTestHelpert(t, notif, 150*time.Millisecond)
	notif.AssertCalled(t, "NotifyStateChange", general.StateHealthy)
}

func TestHealthyThenUnhealthy(t *testing.T) {
	notif := &notifierMock{}
	engineTestHelpert(t, notif, time.Second)

	notif.AssertCalled(t, "NotifyStateChange", general.StateHealthy)
	notif.AssertCalled(t, "NotifyStateChange", general.StateUnhealthy)
}

func engineTestHelpert(t *testing.T, notif *notifierMock, duration time.Duration) {
	notif.On("NotifyStateChange", mock.Anything).Return(nil)
	engine := New(testEngineConfig(), notif)
	err := engine.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(duration))
	defer cancel()
	go engine.Start(ctx)

	err = engine.Ping(somePassphrase)
	assert.NoError(t, err)

	<-ctx.Done()
}

type notifierMock struct {
	mock.Mock
}

func (n *notifierMock) Init() error {
	args := n.Called()

	return args.Error(0)
}

func (n *notifierMock) NotifyStateChange(state general.State) error {
	log.Printf("notify state change %d\n", state)
	args := n.Called(state)

	return args.Error(0)
}
