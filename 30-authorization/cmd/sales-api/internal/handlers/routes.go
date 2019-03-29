package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/ardanlabs/garagesale/internal/mid"
	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/web"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, errorLog, infoLog *log.Logger, authenticator *auth.Authenticator) http.Handler {

	app := web.New(errorLog, mid.RequestLogger(infoLog), web.ErrorHandler(errorLog), mid.Metrics)

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
		p := NewProducts(db, infoLog)

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
