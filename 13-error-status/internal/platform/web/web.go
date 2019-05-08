package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// Handler is the signature used by all application handlers in this service.
type Handler func(http.ResponseWriter, *http.Request) error

// App is the entrypoint into our application and what controls the context of
// each request. Feel free to add any configuration data/logic on this type.
type App struct {
	log *log.Logger
	mux *chi.Mux
}

// NewApp constructs an App to handle a set of routes.
func NewApp(log *log.Logger) *App {
	return &App{
		log: log,
		mux: chi.NewRouter(),
	}
}

// Handle associates a handler function with an HTTP Method and URL pattern.
//
// It converts our custom handler type to the std lib Handler type. It captures
// errors from the handler and serves them to the client in a uniform way.
func (a *App) Handle(method, url string, h Handler) {

	fn := func(w http.ResponseWriter, r *http.Request) {

		// Call the handler and catch any propagated error.
		err := h(w, r)

		if err != nil {

			// If the error was of the type *Error, the handler has
			// a specific status code and error to return. If not, the
			// handler sent any arbitrary error value so use 500.
			webErr, ok := errors.Cause(err).(*Error)
			if !ok {
				webErr = &Error{
					Err:    err,
					Status: http.StatusInternalServerError,
				}
			}

			// Log the error.
			log.Printf("ERROR : %+v", err)

			// Determine the error message service users will see. If the status
			// code is under 500 then it is a "human readable" error that was
			// intended for users to see. If the status code is 500 or higher (the
			// default) then use a generic error message.
			var errStr string
			if webErr.Status < http.StatusInternalServerError {
				errStr = webErr.Err.Error()
			} else {
				errStr = http.StatusText(webErr.Status)
			}

			// Respond with the error type we send to clients.
			res := ErrorResponse{
				Error: errStr,
			}

			Respond(w, res, webErr.Status)
		}
	}

	a.mux.MethodFunc(method, url, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
