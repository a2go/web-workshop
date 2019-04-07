package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Handler is the signature used by all application handlers in this service.
type Handler func(http.ResponseWriter, *http.Request) error

// App is the entrypoint into our application and what controls the context of
// each request. Feel free to add any configuration data/logic on this type.
type App struct {
	log *log.Logger
	mux *chi.Mux
}

// New constructs an App to handle a set of routes.
func New(log *log.Logger) *App {
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
			serr := ToStatusError(err)

			// If the error is an internal issue then log it.
			// Do not log errors that come from client requests.
			if serr.Status >= http.StatusInternalServerError {
				log.Printf("%+v", err)
			}

			// Tell the client about the error.
			res := ErrorResponse{
				Error:  serr.ExternalError(),
				Fields: serr.Fields,
			}

			Respond(w, res, serr.Status)
		}
	}

	a.mux.MethodFunc(method, url, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
