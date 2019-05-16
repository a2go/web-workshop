package handlers

import (
	"log"
	"net/http"

	"github.com/ardanlabs/garagesale/internal/mid"
	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, log *log.Logger, authenticator *auth.Authenticator) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(log, mid.Logger(log), mid.Errors(log), mid.Metrics())

	{
		// Register health check handler. This route is not authenticated.
		c := Checks{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		// Register user handlers.
		u := Users{db: db, authenticator: authenticator}

		// The token route can't be authenticated because they need this route to
		// get the token in the first place.
		app.Handle(http.MethodGet, "/v1/users/token", u.Token)
	}

	{
		// Register Product handlers. Ensure all routes are authenticated.
		p := Products{db: db, log: log}

		app.Handle(http.MethodGet, "/v1/products", p.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/v1/products", p.Create, mid.Authenticate(authenticator))
		app.Handle(http.MethodPut, "/v1/products/{id}", p.Update, mid.Authenticate(authenticator))
		app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete, mid.Authenticate(authenticator))

		app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales, mid.Authenticate(authenticator))
	}

	return app
}
