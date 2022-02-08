package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/idoberko2/home_health_be/general"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errTest = errors.New("some error")

func testScheduler(stateHandler StateHandler) Scheduler {
	return NewScheduler(SchedulerConfig{
		SampleRate: 100 * time.Millisecond,
	}, stateHandler)
}

func TestSchedulerStart(t *testing.T) {
	stateHandlerMock := &stateHandler{}
	stateHandlerMock.On("CheckState").Return(general.StateHealthy, nil)
	stateHandlerMock.On("OnStateCheck", general.StateHealthy).Return(nil)
	scheduler := testScheduler(stateHandlerMock)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(510*time.Millisecond))
	defer cancel()

	errReporter := make(chan error)
	go scheduler.Start(ctx, errReporter)

	select {
	case err := <-errReporter:
		t.Errorf("unexpected error %s", err)
	case <-ctx.Done():
	}

	stateHandlerMock.AssertNumberOfCalls(t, "CheckState", 5)
	stateHandlerMock.AssertNumberOfCalls(t, "OnStateCheck", 5)
}

func TestSchedulerCheckStateError(t *testing.T) {
	stateHandlerMock := &stateHandler{}
	stateHandlerMock.On("CheckState").Return(general.StateHealthy, errTest)
	stateHandlerMock.On("OnStateCheck", general.StateHealthy).Return(nil)
	scheduler := testScheduler(stateHandlerMock)
	schedulerErrHelper(t, scheduler)

	stateHandlerMock.AssertNumberOfCalls(t, "CheckState", 1)
	stateHandlerMock.AssertNotCalled(t, "OnStateCheck", mock.Anything)
}

func TestSchedulerOnStateCheckError(t *testing.T) {
	stateHandlerMock := &stateHandler{}
	stateHandlerMock.On("CheckState").Return(general.StateHealthy, nil)
	stateHandlerMock.On("OnStateCheck", general.StateHealthy).Return(errTest)
	scheduler := testScheduler(stateHandlerMock)
	schedulerErrHelper(t, scheduler)

	stateHandlerMock.AssertNumberOfCalls(t, "CheckState", 1)
	stateHandlerMock.AssertNumberOfCalls(t, "OnStateCheck", 1)
}

func schedulerErrHelper(t *testing.T, scheduler Scheduler) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(110*time.Millisecond))
	defer cancel()

	errReporter := make(chan error)
	go scheduler.Start(ctx, errReporter)

	var expected error
	select {
	case err := <-errReporter:
		expected = err
	case <-ctx.Done():
	}

	assert.Equal(t, expected, errTest)
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
