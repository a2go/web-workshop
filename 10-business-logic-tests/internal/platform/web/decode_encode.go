package web

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/service-training/10-business-logic-tests/internal/platform/log"
)

type HandlerFunc func(r *http.Request) (interface{}, error)

func Decode(r *http.Request, x interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(x); err != nil {
		return ErrorWithStatus(err, http.StatusBadRequest)
	}

	return nil
}

func Encode(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := f(r)
		if err != nil {
			log.Log("request resulted in error", "error", err)

			serr := toStatusError(err)
			w.WriteHeader(serr.status)

			res = struct {
				Error string `json:"error"`
			}{serr.ExternalError()}
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Log("responding with json", "error", err)
		}
	}
}
