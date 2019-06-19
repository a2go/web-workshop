package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("product not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// List gets all Products from the database.
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	products := []Product{}
	const q = `SELECT
			p.*,
			COALESCE(SUM(s.quantity) ,0) AS sold,
			COALESCE(SUM(s.paid), 0) AS revenue
		FROM products AS p
		LEFT JOIN sales AS s ON p.product_id = s.product_id
		GROUP BY p.product_id`

	if err := db.SelectContext(ctx, &products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Create adds a Product to the database. It returns the created Product with
// fields like ID and DateCreated populated..
func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
		INSERT INTO products
		(product_id, name, cost, quantity, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.ExecContext(ctx, q,
		p.ID, p.Name,
		p.Cost, p.Quantity,
		p.DateCreated, p.DateUpdated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}

	return &p, nil
}

// Retrieve finds the product identified by a given ID.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var p Product

	const q = `SELECT
			p.*,
			COALESCE(SUM(s.quantity), 0) AS sold,
			COALESCE(SUM(s.paid), 0) AS revenue
		FROM products AS p
		LEFT JOIN sales AS s ON p.product_id = s.product_id
		WHERE p.product_id = $1
		GROUP BY p.product_id`

	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single product")
	}

	return &p, nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func Update(ctx context.Context, db *sqlx.DB, id string, update UpdateProduct, now time.Time) error {
	p, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Cost != nil {
		p.Cost = *update.Cost
	}
	if update.Quantity != nil {
		p.Quantity = *update.Quantity
	}
	p.DateUpdated = now

	const q = `UPDATE products SET
		"name" = $2,
		"cost" = $3,
		"quantity" = $4,
		"date_updated" = $5
		WHERE product_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		p.Name, p.Cost,
		p.Quantity, p.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}
