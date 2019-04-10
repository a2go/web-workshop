package mid

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// Logger writes some information about the request to the logs in the
// format: (200) GET /foo -> IP ADDR (latency)
func (mw *Middleware) Logger(before web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ctx, span := trace.StartSpan(ctx, "internal.mid.RequestLogger")
		defer span.End()

		v, ok := ctx.Value(web.KeyValues).(*web.Values)
		if !ok {
			return errors.New("web value missing from context")
		}

		err := before(ctx, w, r)

		mw.Log.Printf("(%d) : %s %s -> %s (%s)",
			v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(v.Start),
		)

		// Return the error so it can be handled further up the chain.
		return err
	}

	return h
}
