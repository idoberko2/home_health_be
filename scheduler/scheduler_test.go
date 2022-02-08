package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/idoberko2/home_health_be/general"
	"github.com/stretchr/testify/mock"
)

func testScheduler(stateHandler StateHandler) Scheduler {
	return NewScheduler(SchedulerConfig{
		StateHandler: stateHandler,
		SampleRate:   100 * time.Millisecond,
	})
}

func TestSchedulerStart(t *testing.T) {
	stateHandlerMock := &stateHandler{}
	stateHandlerMock.On("CheckState").Return(general.StateHealthy, nil)
	stateHandlerMock.On("OnStateCheck", general.StateHealthy).Return(nil)
	scheduler := testScheduler(stateHandlerMock)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(510*time.Millisecond))
	defer cancel()

	errReporter := make(chan error)
	scheduler.Start(ctx, errReporter)

	select {
	case err := <-errReporter:
		t.Errorf("unexpected error %s", err)
	case <-ctx.Done():
	}

	stateHandlerMock.AssertNumberOfCalls(t, "CheckState", 5)
	stateHandlerMock.AssertNumberOfCalls(t, "OnStateCheck", 5)
}

type stateHandler struct {
	mock.Mock
}

func (s *stateHandler) CheckState() (general.State, error) {
	args := s.Called()

	var arg0 general.State
	if args.Get(0) != nil {
		arg0 = args.Get(0).(general.State)
	}

	return arg0, args.Error(1)
}

func (s *stateHandler) OnStateCheck(state general.State) error {
	args := s.Called(state)

	return args.Error(0)
}
