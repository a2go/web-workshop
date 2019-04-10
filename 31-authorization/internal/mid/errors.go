package mid

import (
	"context"
	"net/http"

	"github.com/ardanlabs/garagesale/internal/platform/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func (mw *Middleware) Errors(before web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		// Run the handler chain and catch any propagated error.
		if err := before(ctx, w, r); err != nil {
			serr := web.ToStatusError(err)

			// If the error is an internal issue then log it.
			// Do not log errors that come from client requests.
			if serr.Status >= http.StatusInternalServerError {
				mw.Log.Printf("%+v", err)
			}

			// Tell the client about the error.
			res := web.ErrorResponse{
				Error:  serr.ExternalError(),
				Fields: serr.Fields,
			}

			if err := web.Respond(ctx, w, res, serr.Status); err != nil {
				return err
			}
		}

		// Return nil to indicate the error has been handled.
		return nil
	}

	return h
}
