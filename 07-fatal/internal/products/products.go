package products

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Product is an item we sell.
type Product struct {
	ID       string `db:"product_id"`
	Name     string `db:"name"`
	Cost     int    `db:"cost"`
	Quantity int    `db:"quantity"`
}

// List gets all Products from the database.
func List(db *sqlx.DB) ([]Product, error) {
	var products []Product

	if err := db.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}
