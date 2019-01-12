package web

import (
	"encoding/json"
	"net/http"

	"github.com/ardanlabs/service-training/10-business-logic-tests/internal/platform/log"
)

// HandlerFunc is the signature used by all application handlers in this
// service. It should return the value to encode to the client. Any error
// returned by this handler will be encoded to the client with a custom status
// code.
type HandlerFunc func(r *http.Request) (interface{}, error)

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return ErrorWithStatus(err, http.StatusBadRequest)
	}

	return nil
}

// Encode wraps a HandlerFunc as defined in this package in a new function
// compatible with net/http. The new function calls the provided function and
// knows how to encode the returned value for the client response.
func Encode(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := fn(r)
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
