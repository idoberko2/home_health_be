package main

import "time"

type SchedulerConfig struct {
	OnStateChange func(State) error
	SampleRate    time.Duration
}
