package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Handler is the signature used by all application handlers in a service.
type Handler func(http.ResponseWriter, *http.Request)

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
func (a *App) Handle(method, url string, h Handler) {
	a.mux.Method(method, url, http.HandlerFunc(h))
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
