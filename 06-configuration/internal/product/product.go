package product

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// List gets all Products from the database.
func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT * FROM products`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}
