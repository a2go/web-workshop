package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
)

// Open knows how to open a database connection.
func Open() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
