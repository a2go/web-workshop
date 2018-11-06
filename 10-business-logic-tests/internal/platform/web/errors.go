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
	return se.Error()
}

func statusFromError(err error) int {
	if se, ok := err.(*statusError); ok {
		return se.status
	}
	return http.StatusInternalServerError
}
