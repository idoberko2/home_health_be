package scheduler

import (
	"time"
)

type SchedulerConfig struct {
	StateHandler
	SampleRate time.Duration
}
