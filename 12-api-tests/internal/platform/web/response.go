package web

import (
	"encoding/json"
	"net/http"
)

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
