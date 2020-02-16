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

	p := Products{db: db, log: log}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)

	return app
}
