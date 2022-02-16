package app

import (
	"context"
	"os"

	"github.com/idoberko2/home_health_be/engine"
	"github.com/idoberko2/home_health_be/notifier"

	"github.com/joho/godotenv"
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
	avoiddotenv := os.Getenv("HC_AVOID_DOTENV")
	if avoiddotenv == "" {
		err := godotenv.Load()
		if err != nil {
			log.WithError(err).Fatal("Error loading .env file")
		}
	}

	ctx := context.Background()
	g, grpCtx := errgroup.WithContext(ctx)
	g.Go(getStartEngine(grpCtx))

	if err := g.Wait(); err != nil {
		log.WithError(err).Error("did not finish properly")
	}
}

func getStartEngine(ctx context.Context) func() error {
	return func() error {
		cfg, errCfg := ReadEngineConfig()
		if errCfg != nil {
			return errCfg
		}

		engine := engine.New(cfg, notifier.NewLogNotifier())
		if err := engine.Init(); err != nil {
			return err
		}

		return engine.Start(ctx)
	}
}
