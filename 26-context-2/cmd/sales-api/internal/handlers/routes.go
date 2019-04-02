package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/ardanlabs/garagesale/internal/mid"
	"github.com/ardanlabs/garagesale/internal/platform/web"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, errorLog, infoLog *log.Logger) http.Handler {

	app := web.New(errorLog, mid.RequestLogger(infoLog), web.ErrorHandler(errorLog), mid.Metrics)

	{
		c := Checks{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		p := Products{db: db, log: infoLog}

		app.Handle(http.MethodGet, "/v1/products", p.List)
		app.Handle(http.MethodGet, "/v1/products/{id}", p.Get)
		app.Handle(http.MethodPost, "/v1/products", p.Create)
		app.Handle(http.MethodPut, "/v1/products/{id}", p.Update)
		app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete)

		app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
		app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)
	}

	return app
}
