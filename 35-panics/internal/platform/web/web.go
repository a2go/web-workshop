package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	TraceID    string
	StatusCode int
	Start      time.Time
}

// Handler is the signature used by all application handlers in this service.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

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
func (a *App) Handle(method, url string, h Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	h = wrapMiddleware(mw, h)

	// Add the application's general middleware to the handler chain.
	h = wrapMiddleware(a.mw, h)

	// Create a function that conforms to the std lib defintioin of a handler.
	// This is the first thing that will be executed when this route is called.
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "internal.platform.web")
		defer span.End()

		// Create a Values struct to record state for the request. Store the
		// address in the request's context so it is sent down the call chain.
		v := Values{
			TraceID: span.SpanContext().TraceID.String(),
			Start:   time.Now(),
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		// Run the handler chain and catch any propagated error.
		if err := h(ctx, w, r); err != nil {
			log.Printf("Unhandled error: %+v", err)
		}
	}

	// Create an OpenCensus HTTP Handler which wraps our chain. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if an client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/
	och := &ochttp.Handler{
		Handler:     http.HandlerFunc(fn),
		Propagation: &tracecontext.HTTPFormat{},
	}

	// Register the
	a.mux.Method(method, url, och)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
