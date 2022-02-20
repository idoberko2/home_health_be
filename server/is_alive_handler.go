package server

import (
	"net/http"

	"github.com/idoberko2/home_health_be/engine"
	log "github.com/sirupsen/logrus"
)

func GetIsAliveHandler(e engine.Engine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Alive!")
		w.WriteHeader(http.StatusOK)
	}
}
