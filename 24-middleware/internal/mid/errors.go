package mid

import (
	"net/http"

	"github.com/ardanlabs/garagesale/internal/platform/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func (mw *Middleware) Errors(before web.Handler) web.Handler {
	h := func(w http.ResponseWriter, r *http.Request) error {

		// Run the handler chain and catch any propagated error.
		if err := before(w, r); err != nil {

			// Convert the error interface variable to the concrete type
			// *web.StatusError to find the appropriate HTTP status.
			serr := web.NewStatusError(err)

			// If the error is an internal issue then log the error message.
			// Do not log error messages that come from client requests.
			if serr.Status >= http.StatusInternalServerError {
				mw.Log.Printf("%+v", err)
			}

			// Respond with the error type we send to clients.
			res := web.ErrorResponse{
				Error:  serr.String(),
				Fields: serr.Fields,
			}

			if err := web.Respond(w, res, serr.Status); err != nil {
				return err
			}
		}

		// Return nil to indicate the error has been handled.
		return nil
	}

	return h
}
