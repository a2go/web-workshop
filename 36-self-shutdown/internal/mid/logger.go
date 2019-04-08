package mid

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// RequestLogger writes some information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
type RequestLogger struct {
	Log *log.Logger
}

// Handle is the actual Middleware function to be ran when building the chain.
func (mw *RequestLogger) Handle(before web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ctx, span := trace.StartSpan(ctx, "internal.mid.RequestLogger")
		defer span.End()

		// If the context is missing this value, request the service
		// to be shutdown gracefully.
		v, ok := ctx.Value(web.KeyValues).(*web.Values)
		if !ok {
			return web.Shutdown("web value missing from context")
		}

		err := before(ctx, w, r)

		mw.Log.Printf("%s : (%d) : %s %s -> %s (%s)",
			v.TraceID, v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(v.Start),
		)

		// Return the error so it can be handled further up the chain.
		return err
	}

	return h
}
