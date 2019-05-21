package web

import (
	"encoding/json"
	"net/http"
)

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	return nil
}
