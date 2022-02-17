package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/idoberko2/home_health_be/engine"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	Start(ctx context.Context) error
}

func New(e engine.Engine, cfg ServerConfig) Server {
	r := mux.NewRouter()
	r.HandleFunc("/ping", GetPingHandler(e)).Methods(http.MethodPost)

	return &server{
		r:   r,
		cfg: cfg,
	}
}

type server struct {
	r   *mux.Router
	cfg ServerConfig
}

func (s *server) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		WriteTimeout: s.cfg.WriteTimeout,
		ReadTimeout:  s.cfg.ReadTimeout,
		IdleTimeout:  s.cfg.IdleTimeout,
		Handler:      s.r,
	}

	errChan := make(chan error)
	go func() {
		log.Info("starting server...")
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ServerShutdownTimeout)
	defer cancel()

	log.Info("shutting down server...")
	return srv.Shutdown(ctx)
}
