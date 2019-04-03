package tests

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/database/databasetest"
	"github.com/ardanlabs/garagesale/internal/platform/database/schema"
)

// Test owns state for running and shutting down tests.
type Test struct {
	DB            *sqlx.DB
	Log           *log.Logger
	Authenticator *auth.Authenticator

	cleanup func()
}

// New creates a database, seeds it, constructs an authenticator.
func New(t *testing.T) *Test {
	t.Helper()

	var test Test

	// Initialize and seed database. Store the cleanup function call later.
	test.DB, test.cleanup = databasetest.Setup(t)

	if err := schema.Seed(test.DB.DB); err != nil {
		t.Fatal(err)
	}

	// Create the logger to use.
	test.Log = log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// Create RSA keys to enable authentication in our service.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Build an authenticator using this static key.
	kid := "4754d86b-7a6d-4df5-9c65-224741361492"
	kf := auth.NewSingleKeyFunc(kid, key.Public().(*rsa.PublicKey))
	test.Authenticator, err = auth.NewAuthenticator(key, kid, "RS256", kf)
	if err != nil {
		t.Fatal(err)
	}

	return &test
}

// Teardown releases any resources used for the test.
func (test *Test) Teardown() {
	test.cleanup()
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}
