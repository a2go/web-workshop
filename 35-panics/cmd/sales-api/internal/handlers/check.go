package handlers

import (
	"context"
	"net/http"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"
)

// Checks holds handlers for service orchestration checks.
type Checks struct {
	db *sqlx.DB
}

// Health returns a 200 OK status if the service is ready to
// receive requests.
func (s *Checks) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Checks.Health")
	defer span.End()

	var response struct {
		Status string `json:"status"`
	}

	// Check if the database is ready.
	if err := s.db.PingContext(ctx); err != nil {

		// If the database is not ready we will tell the client and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		response.Status = "not ready"
		return web.Respond(ctx, w, response, http.StatusInternalServerError)
	}

	response.Status = "ok"
	return web.Respond(ctx, w, response, http.StatusOK)
}
