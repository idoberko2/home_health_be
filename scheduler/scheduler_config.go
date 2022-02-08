package scheduler

import (
	"time"

	"github.com/idoberko2/home_health_be/general"
)

type SchedulerConfig struct {
	OnStateChange func(general.State) error
	SampleRate    time.Duration
}
