package main

type Scheduler interface {
	Start() error
	Stop() error
}

func NewScheduler(cfg SchedulerConfig) Scheduler {
	return &scheduler{
		cfg: cfg,
	}
}

type scheduler struct {
	cfg SchedulerConfig
}

func (s *scheduler) Start() error {
	return nil
}

func (s *scheduler) Stop() error {
	return nil
}
