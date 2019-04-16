package mid

import (
	"context"
	"net/http"
	"strings"

	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/pkg/errors"
)

// Authenticate validates a JWT from the `Authorization` header.
func (mw *Middleware) Authenticate(after web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		authHdr := r.Header.Get("Authorization")
		if authHdr == "" {
			err := errors.New("missing Authorization header")
			return web.WrapErrorWithStatus(err, http.StatusUnauthorized)
		}

		tknStr, err := parseAuthHeader(authHdr)
		if err != nil {
			return web.WrapErrorWithStatus(err, http.StatusUnauthorized)
		}

		claims, err := mw.Authenticator.ParseClaims(tknStr)
		if err != nil {
			return web.WrapErrorWithStatus(err, http.StatusUnauthorized)
		}

		// Add claims to the context so they can be retrieved later.
		ctx = context.WithValue(ctx, auth.Key, claims)

		return after(ctx, w, r)
	}

	return h
}

// parseAuthHeader parses an authorization header. Expected header is of
// the format `Bearer <token>`.
func parseAuthHeader(bearerStr string) (string, error) {
	split := strings.Split(bearerStr, " ")
	if len(split) != 2 || strings.ToLower(split[0]) != "bearer" {
		return "", errors.New("Expected Authorization header format: Bearer <token>")
	}

	return split[1], nil
}
