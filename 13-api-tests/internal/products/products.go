package products

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("product not found")

	// ErrInvalidID is used when a specific Product is requested but does not exist.
	ErrInvalidID = errors.New("ID is not in its proper form")
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
	var products []Product

	if err := db.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Create uses the provided *Product to insert a new product record. The ID
// field provided is populated.
func Create(db *sqlx.DB, p *Product) error {
	p.ID = uuid.New().String()

	_, err := db.Exec(`
		INSERT INTO products
		(product_id, name, cost, quantity)
		VALUES ($1, $2, $3, $4)`,
		p.ID, p.Name, p.Cost, p.Quantity,
	)
	if err != nil {
		return errors.Wrap(err, "inserting product")
	}

	return nil
}

// Get finds the product identified by a given ID.
func Get(db *sqlx.DB, id string) (*Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var p Product

	err := db.Get(&p, `
		SELECT * FROM products
		WHERE product_id = $1`,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single product")
	}

	return &p, nil
}
