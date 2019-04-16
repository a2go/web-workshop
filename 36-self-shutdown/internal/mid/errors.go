package mid

import (
	"context"
	"net/http"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func (mw *Middleware) Errors(before web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ctx, span := trace.StartSpan(ctx, "internal.mid.ErrorHandler")
		defer span.End()

		// If the context is missing this value, request the service
		// to be shutdown gracefully.
		v, ok := ctx.Value(web.KeyValues).(*web.Values)
		if !ok {
			return web.Shutdown("web value missing from context")
		}

		// Run the handler chain and catch any propagated error.
		if err := before(ctx, w, r); err != nil {

			// Convert the error interface variable to the concrete type
			// *web.StatusError to find the appropriate HTTP status.
			serr := web.NewStatusError(err)

			// If the error is an internal issue then log the error message.
			// Do not log error messages that come from client requests.
			if serr.Status >= http.StatusInternalServerError {
				mw.Log.Printf("%s : %+v", v.TraceID, err)
			}

			// Respond with the error type we send to clients.
			res := web.ErrorResponse{
				Error:  serr.String(),
				Fields: serr.Fields,
			}

			if err := web.Respond(ctx, w, res, serr.Status); err != nil {
				return err
			}

			// If the error that was just handled was the special Shutdown error
			// then let that return.
			if web.IsShutdown(err) {
				return err
			}
		}

		// Return nil to indicate the error has been handled.
		return nil
	}

	return h
}
