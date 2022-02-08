package scheduler

import (
	"context"
	"time"
)

type Scheduler interface {
	Start(ctx context.Context, errReporter chan error)
}

func NewScheduler(cfg SchedulerConfig) Scheduler {
	return &scheduler{
		cfg: cfg,
	}
}

type scheduler struct {
	cfg SchedulerConfig
}

func (s *scheduler) Start(ctx context.Context, errReporter chan error) {
	ticker := time.NewTicker(s.cfg.SampleRate)
	for {
		select {
		case <-ticker.C:
			{
				state, err := s.cfg.StateHandler.CheckState()
				if err != nil {
					errReporter <- err
					break
				}
				if err := s.cfg.StateHandler.OnStateCheck(state); err != nil {
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
