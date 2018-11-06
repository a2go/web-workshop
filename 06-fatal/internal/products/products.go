package products

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Product struct {
	ID       int    `db:"product_id"`
	Name     string `db:"name"`
	Cost     int    `db:"cost"`
	Quantity int    `db:"quantity"`
}

func DBConn(user, pass, host, name string, disableSSL bool) string {
	sslMode := "require"
	if disableSSL {
		sslMode = "disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s&timezone=utc",
		user, pass, host, name, sslMode)
}

func List(db *sqlx.DB) ([]Product, error) {
	var products []Product

	// TODO: Talk about issues of using '*' after talking about migrations.
	if err := db.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}
