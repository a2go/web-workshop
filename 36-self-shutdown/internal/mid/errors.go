package mid

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// ErrorHandler holds a Middlware that handles errors coming out of the call
// chain. It detects normal applications errors which are used to respond to
// the client in a uniform way. Unexpected errors (status >= 500) are logged.
type ErrorHandler struct {
	Log *log.Logger
}

// Handle is the actual Middleware function to be ran when building the chain.
func (mw *ErrorHandler) Handle(before web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ctx, span := trace.StartSpan(ctx, "internal.mid.ErrorHandler")
		defer span.End()

		v, ok := ctx.Value(web.KeyValues).(*web.Values)
		if !ok {
			return errors.New("web value missing from context")
		}

		// Run the handler chain and catch any propagated error.
		if err := before(ctx, w, r); err != nil {
			serr := web.ToStatusError(err)

			// If the error is an internal issue then log it.
			// Do not log errors that come from client requests.
			if serr.Status >= http.StatusInternalServerError {
				mw.Log.Printf("%s : %+v", v.TraceID, err)
			}

			// Tell the client about the error.
			res := web.ErrorResponse{
				Error:  serr.ExternalError(),
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
