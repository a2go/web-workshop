package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/ardanlabs/garagesale/internal/mid"
	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(shutdown chan os.Signal, db *sqlx.DB, log *log.Logger, authenticator *auth.Authenticator) http.Handler {

	// Create the middleware that log requests and handle errors.
	loggerMW := mid.RequestLogger{Log: log}
	errorMW := mid.ErrorHandler{Log: log}

	// Construct the web.App which holds all routes as well as a chain of common Middleware.
	app := web.New(shutdown, log, loggerMW.Handle, errorMW.Handle, mid.Metrics, mid.PanicHandler)

	// Create the middleware that can authenticate and authorize requests.
	authmw := mid.Auth{
		Authenticator: authenticator,
	}

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

		app.Handle(http.MethodGet, "/v1/products", p.List, authmw.Authenticate)
		app.Handle(http.MethodGet, "/v1/products/{id}", p.Get, authmw.Authenticate)
		app.Handle(http.MethodPost, "/v1/products", p.Create, authmw.Authenticate)
		app.Handle(http.MethodPut, "/v1/products/{id}", p.Update, authmw.Authenticate)
		app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete, authmw.Authenticate, authmw.HasRole(auth.RoleAdmin))

		app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale, authmw.Authenticate, authmw.HasRole(auth.RoleAdmin))
		app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales, authmw.Authenticate)
	}

	return app
}
