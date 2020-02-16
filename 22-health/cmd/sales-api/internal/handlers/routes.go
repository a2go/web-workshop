package handlers

import (
	"log"
	"net/http"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, log *log.Logger) http.Handler {

	app := web.NewApp(log)

	{
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		p := Products{db: db, log: log}

		app.Handle(http.MethodGet, "/v1/products", p.List)
		app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
		app.Handle(http.MethodPost, "/v1/products", p.Create)
		app.Handle(http.MethodPut, "/v1/products/{id}", p.Update)
		app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete)

		app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
		app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)
	}

	return app
}
