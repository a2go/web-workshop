package web

import (
	"encoding/json"
	"net/http"
)

// Handler is the signature used by all application handlers in this service.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Run is the entry point for all handlers. It converts our custom handler type
// to the std lib Handler type. It captures errors returned from the handler
// and serves them to the client in a uniform way.
func Run(h Handler) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {

		// Call the handler and catch any propagated error.
		err := h(w, r)

		if err != nil {
			serr := toStatusError(err)

			res := struct {
				Error  string       `json:"error"`
				Fields []fieldError `json:"fields,omitempty"`
			}{
				Error:  serr.ExternalError(),
				Fields: serr.fields,
			}

			Encode(w, res, serr.status)
		}
	}

	return fn
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return ErrorWithStatus(err, http.StatusBadRequest)
	}

	return validateFields(val)
}

// Encode converts a Go value to JSON and sends it to the client.
func Encode(w http.ResponseWriter, data interface{}, status int) error {

	// Convert the response value to JSON.
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Respond with the provided JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(res)); err != nil {
		return err
	}

	return nil
}
