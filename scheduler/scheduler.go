package scheduler

import "context"

type Scheduler interface {
	Start(ctx context.Context) error
}

func NewScheduler(cfg SchedulerConfig) Scheduler {
	return &scheduler{
		cfg: cfg,
	}
}

type scheduler struct {
	cfg SchedulerConfig
}

func (s *scheduler) Start(ctx context.Context) error {
	return nil
}
