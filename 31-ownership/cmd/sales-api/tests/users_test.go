package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/garagesale/cmd/sales-api/internal/handlers"
	"github.com/ardanlabs/garagesale/internal/platform/tests"
)

// TestUsers runs a series of tests to exercise User behavior.
func TestUsers(t *testing.T) {
	test := tests.New(t)
	defer test.Teardown()

	ut := UserTests{app: handlers.API(test.DB, test.Log, test.Authenticator)}

	t.Run("TokenRequireAuth", ut.TokenRequireAuth)
	t.Run("TokenDenyUnknown", ut.TokenDenyUnknown)
	t.Run("TokenDenyBadPassword", ut.TokenDenyBadPassword)
	t.Run("TokenSuccess", ut.TokenSuccess)
}

// UserTests holds methods for each user subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type UserTests struct {
	app http.Handler
}

// TokenRequireAuth ensures that requests with no authentication are denied.
func (ut *UserTests) TokenRequireAuth(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/users/token", nil)
	resp := httptest.NewRecorder()

	ut.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusUnauthorized, resp.Code)
	}
}

// TokenDenyUnknown ensures that users with an unrecognized email aren't given a token.
func (ut *UserTests) TokenDenyUnknown(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/users/token", nil)
	resp := httptest.NewRecorder()

	req.SetBasicAuth("unknown@example.com", "gophers")

	ut.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusUnauthorized, resp.Code)
	}
}

// TokenDenyBadPassword ensures that a known user with a bad password is not authenticated.
func (ut *UserTests) TokenDenyBadPassword(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/users/token", nil)
	resp := httptest.NewRecorder()

	req.SetBasicAuth("admin@example.com", "GOPHERS")

	ut.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusUnauthorized, resp.Code)
	}
}

// TokenSuccess tests that a known user with a good password gets a token.
func (ut *UserTests) TokenSuccess(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/users/token", nil)
	resp := httptest.NewRecorder()

	req.SetBasicAuth("admin@example.com", "gophers")

	ut.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var got map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	if len(got) != 1 {
		t.Error("unexpected values in token response")
	}

	if got["token"] == "" {
		t.Fatal("token was not in response")
	}
}
