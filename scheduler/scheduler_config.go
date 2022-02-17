package scheduler

import (
	"time"
)

type SchedulerConfig struct {
	SampleRate time.Duration `split_words:"true" default:"10s"`
}
