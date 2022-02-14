package app

import (
	"context"

	"github.com/idoberko2/home_health_be/engine"
	"github.com/idoberko2/home_health_be/notifier"

	log "github.com/sirupsen/logrus"
)

type App interface {
	Run()
}

func New() App {
	return &app{}
}

type app struct{}

func (a *app) Run() {
	ctx := context.Background()

	engine := engine.New(engine.EngineConfig{}, notifier.NewLogNotifier())
	if err := engine.Init(); err != nil {
		log.WithError(err).Fatal("error initializing engine")
	}

	if err := engine.Start(ctx); err != nil {
		log.WithError(err).Fatal("error starting engine")
	}
}
