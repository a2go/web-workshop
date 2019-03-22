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
	mw  []Middleware
}

// New constructs an App to handle a set of routes. Any Middleware provided
// will be ran for every request.
func New(log *log.Logger, mw ...Middleware) *App {
	return &App{
		log: log,
		mux: chi.NewRouter(),
		mw:  mw,
	}
}

// Handle associates a handler function with an HTTP Method and URL pattern.
//
// It converts our custom handler type to the std lib Handler type. It captures
// errors from the handler and serves them to the client in a uniform way.
func (a *App) Handle(method, url string, h Handler) {

	// wrap the provided handler in the application's middleware.
	h = wrapMiddleware(h, a.mw)

	fn := func(w http.ResponseWriter, r *http.Request) {

		// Run the handler chain and catch any propagated error.
		err := h(w, r)

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

			Encode(w, res, serr.status)
		}
	}

	a.mux.MethodFunc(method, url, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
