package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	appConfig, err := ReadAppConfig()
	if err != nil {
		log.WithError(err).Fatal("error reading app config")
	}

	if !appConfig.AvoidDotEnv {
		err := godotenv.Load()
		if err != nil {
			log.WithError(err).Fatal("Error loading .env file")
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdown(cancel)

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

func gracefulShutdown(terminate func()) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		s := <-sigc
		log.WithField("signal", s.String()).Info("received signal. terminating...")
		terminate()
	}()
}
