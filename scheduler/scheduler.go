package scheduler

import (
	"context"
	"time"
)

type Scheduler interface {
	Start(ctx context.Context, errReporter chan error)
}

func New(cfg SchedulerConfig, stateHandler StateHandler) Scheduler {
	return &scheduler{
		cfg:          cfg,
		stateHandler: stateHandler,
	}
}

type scheduler struct {
	cfg          SchedulerConfig
	stateHandler StateHandler
}

func (s *scheduler) Start(ctx context.Context, errReporter chan error) {
	ticker := time.NewTicker(s.cfg.SampleRate)
	for {
		select {
		case <-ticker.C:
			{
				state, err := s.stateHandler.CheckState()
				if err != nil {
					errReporter <- err
					break
				}
				if err := s.stateHandler.OnStateCheck(state); err != nil {
					errReporter <- err
				}
			}
		case <-ctx.Done():
			{
				return
			}
		}
	}
}
