package server

import (
	"net/http"

	"github.com/idoberko2/home_health_be/engine"
	jparser "github.com/idoberko2/json-request-parser"
)

func GetPingHandler(e engine.Engine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload pingRequestPayload
		if ok := jparser.ParseJSONRequest(w, r, &payload); !ok {
			return
		}

		if err := e.Ping(payload.Passphrase); err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type pingRequestPayload struct {
	Passphrase string
}
