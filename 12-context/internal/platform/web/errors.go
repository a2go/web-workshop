package web

import (
	"net/http"
)

func ErrorWithStatus(err error, status int) error {
	return &statusError{err, status}
}

type statusError struct {
	err    error
	status int
}

func (se *statusError) Error() string {
	return se.err.Error()
}

func (se *statusError) ExternalError() string {
	if se.status < 500 {
		return se.err.Error()
	}
	return http.StatusText(se.status)
}

func toStatusError(err error) *statusError {
	if se, ok := err.(*statusError); ok {
		return se
	}
	return &statusError{err, http.StatusInternalServerError}
}
