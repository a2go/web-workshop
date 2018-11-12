package products

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Product struct {
	ID       int    `db:"product_id"`
	Name     string `db:"name"`
	Cost     int    `db:"cost"`
	Quantity int    `db:"quantity"`
}

func List(db *sqlx.DB) ([]Product, error) {
	var products []Product

	// TODO: Talk about issues of using '*' after talking about migrations.
	if err := db.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}
