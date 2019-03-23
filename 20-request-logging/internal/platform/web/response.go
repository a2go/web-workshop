package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// Respond converts a Go value to JSON and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, status int) error {

	// Set the status code for the request logger middleware.
	v := ctx.Value(KeyValues).(*Values)
	v.StatusCode = status

	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return nil
	}

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
