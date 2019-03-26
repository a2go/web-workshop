package web

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// errorResponse is the form used for API responses from failures in the API.
type errorResponse struct {
	Error  string       `json:"error"`
	Fields []fieldError `json:"fields,omitempty"`
}

// fieldError is used to indicate an error with a specific request field.
type fieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// statusError is used to pass errors through the application with web specific
// context.
type statusError struct {
	err    error
	status int
	fields []fieldError
}

// ErrorWithStatus wraps a provided error with an HTTP status code.
func ErrorWithStatus(err error, status int) error {
	return &statusError{err, status, nil}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (se *statusError) Error() string {
	return se.err.Error()
}

// ExternalError provides "human readable" error messages that are intended for
// service users to see. If the status code is 500 or higher (the default) then
// a generic error message is returned.
//
// The idea is that a developer who creates an error like this intends to let
// the API consumer know the product was not found:
//	ErrorWithStatus(errors.New("product not found"), 404)
//
// However a more serious error like a database failure might include
// information that is not safe to show to API consumers.
func (se *statusError) ExternalError() string {
	if se.status < http.StatusInternalServerError {
		return se.err.Error()
	}
	return http.StatusText(se.status)
}

// toStatusError takes a regular error and converts it to a statusError. If the
// original error is already a *statusError it is returned directly. If not
// then it is defaulted to an error with a 500 status.
func toStatusError(err error) *statusError {
	if se, ok := errors.Cause(err).(*statusError); ok {
		return se
	}
	return &statusError{err, http.StatusInternalServerError, nil}
}

// ErrorHandler creates a middlware that handles errors come out of the call
// chain. It detects normal applications errors which are used to respond to
// the client in a uniform way. Unexpected errors (status >= 500) are logged.
func ErrorHandler(log *log.Logger) Middleware {
	mw := func(before Handler) Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {

			// Run the handler chain and catch any propagated error.
			err := before(w, r)

			if err != nil {
				serr := toStatusError(err)

				// If the error is an internal issue then log it.
				// Do not log errors that come from client requests.
				if serr.status >= http.StatusInternalServerError {
					log.Printf("%+v", err)
				}

				// Tell the client about the error.
				res := errorResponse{
					Error:  serr.ExternalError(),
					Fields: serr.fields,
				}

				if err := Respond(w, res, serr.status); err != nil {
					return err
				}
			}

			// Return nil to indicate the error has been handled.
			return nil
		}
		return h
	}
	return mw
}
