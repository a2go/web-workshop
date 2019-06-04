package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// These are the expected values for Claims.Roles.
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

// NewClaims constructs a Claims value for the identified user. The Claims
// expire within a specified duration of the provided time. Additional fields
// of the Claims can be set after calling NewClaims is desired.
func NewClaims(subject string, roles []string, now time.Time, expires time.Duration) Claims {
	c := Claims{
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(expires).Unix(),
		},
	}

	return c
}
