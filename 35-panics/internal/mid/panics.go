package mid

import (
	"context"
	"net/http"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// PanicHandler recovers from panics and converts the panic to an error so it
// is reported in Metrics and handled in ErrorHandler.
func PanicHandler(next web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		ctx, span := trace.StartSpan(ctx, "internal.mid.PanicHandler")
		defer span.End()

		// Defer a function to recover from a panic and set the err return variable
		// after the fact. Using the errors package will generate a stack trace.
		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("panic: %v", r)
			}
		}()

		// Call the next Handler and set its return value in the err variable.
		return next(ctx, w, r)
	}
	return h
}
