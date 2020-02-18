package handlers

import (
	"net/http"

	"github.com/a2go/garagesale/internal/platform/database"
	"github.com/a2go/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// Check provides support for orchestration health checks.
type Check struct {
	db *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// Health validates the service is healthy and ready to accept requests.
func (c *Check) Health(w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}

	// Check if the database is ready.
	if err := database.StatusCheck(r.Context(), c.db); err != nil {

		// If the database is not ready we will tell the client and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "db not ready"
		return web.Respond(w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return web.Respond(w, health, http.StatusOK)
}
