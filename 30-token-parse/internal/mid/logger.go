package mid

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ardanlabs/garagesale/internal/platform/web"
)

// RequestLogger writes some information about the request to the logs in
// the format: (200) GET /foo -> IP ADDR (latency)
func RequestLogger(log *log.Logger) web.Middleware {

	// Wrap this handler around the next one provided.
	mw := func(before web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web value missing from context")
			}

			err := before(ctx, w, r)

			log.Printf("(%d) : %s %s -> %s (%s)",
				v.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(v.Start),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}
		return h
	}

	return mw
}
