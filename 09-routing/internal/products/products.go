package products

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Product is an item we sell.
type Product struct {
	ID       string `db:"product_id" json:"id"`
	Name     string `db:"name" json:"name"`
	Cost     int    `db:"cost" json:"cost"`
	Quantity int    `db:"quantity" json:"quantity"`
}

// List gets all Products from the database.
func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT * FROM products`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Create uses the provided *Product to insert a new product record. The ID
// field provided is populated.
func Create(db *sqlx.DB, p *Product) error {
	p.ID = uuid.New().String()

	const q = `INSERT INTO products
		(product_id, name, cost, quantity)
		VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(q, p.ID, p.Name, p.Cost, p.Quantity)
	if err != nil {
		return errors.Wrap(err, "inserting product")
	}

	return nil
}

// Get finds the product identified by a given ID.
func Get(db *sqlx.DB, id string) (*Product, error) {
	var p Product

	const q = `SELECT * FROM products WHERE product_id = $1`

	if err := db.Get(&p, q, id); err != nil {
		return nil, errors.Wrap(err, "selecting single product")
	}

	return &p, nil
}
