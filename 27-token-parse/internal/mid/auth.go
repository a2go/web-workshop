package mid

import (
	"context"
	"net/http"
	"strings"

	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/pkg/errors"
)

// Auth is used to authenticate and authorize HTTP requests.
type Auth struct {
	Authenticator *auth.Authenticator
}

// Authenticate validates a JWT from the `Authorization` header.
func (a *Auth) Authenticate(after web.Handler) web.Handler {

	// Wrap this handler around the next one provided.
	h := func(w http.ResponseWriter, r *http.Request) error {
		authHdr := r.Header.Get("Authorization")
		if authHdr == "" {
			err := errors.New("missing Authorization header")
			return web.ErrorWithStatus(err, http.StatusUnauthorized)
		}

		tknStr, err := parseAuthHeader(authHdr)
		if err != nil {
			return web.ErrorWithStatus(err, http.StatusUnauthorized)
		}

		claims, err := a.Authenticator.ParseClaims(tknStr)
		if err != nil {
			return web.ErrorWithStatus(err, http.StatusUnauthorized)
		}

		// Add claims to the context so they can be retrieved later.
		ctx := context.WithValue(r.Context(), auth.Key, claims)
		r = r.WithContext(ctx)

		return after(w, r)
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
