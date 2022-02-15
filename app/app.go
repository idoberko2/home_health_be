package app

import (
	"context"

	"github.com/idoberko2/home_health_be/engine"
	"github.com/idoberko2/home_health_be/notifier"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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

	g, grpCtx := errgroup.WithContext(ctx)
	g.Go(getStartEngine(grpCtx))

	if err := g.Wait(); err != nil {
		log.WithError(err).Error("did not finish properly")
	}
}

func getStartEngine(ctx context.Context) func() error {
	return func() error {
		engine := engine.New(engine.EngineConfig{}, notifier.NewLogNotifier())
		if err := engine.Init(); err != nil {
			return err
		}

		return engine.Start(ctx)
	}
}
