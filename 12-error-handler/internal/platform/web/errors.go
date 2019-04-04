package web

// errorResponse is the form used for API responses from failures in the API.
type errorResponse struct {
	Error string `json:"error"`
}
