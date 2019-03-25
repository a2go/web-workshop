package handlers

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/ardanlabs/garagesale/internal/users"
)

// Users holds handlers for dealing with users.
type Users struct {
	db            *sqlx.DB
	authenticator *auth.Authenticator
}

// Token generates a authentication token for a user. The client must include
// an email and password for the request using HTTP Basic Auth. The user will
// be identified by email and authenticated by their password.
func (s *Users) Token(w http.ResponseWriter, r *http.Request) error {

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return web.ErrorWithStatus(err, http.StatusUnauthorized)
	}

	claims, err := users.Authenticate(r.Context(), s.db, time.Now(), email, pass)
	if err != nil {
		if err == users.ErrAuthenticationFailure {
			return web.ErrorWithStatus(err, http.StatusUnauthorized)
		}
		return errors.Wrap(err, "authenticating user")
	}

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = s.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	return web.Respond(r.Context(), w, tkn, http.StatusOK)
}
