package products

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
)

// Product is an item we sell.
type Product struct {
	ID          string    `db:"product_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewProduct defines the information we need when adding a Product to our
// offerings.
type NewProduct struct {
	Name     string `json:"name" validate:"required"`
	Cost     int    `json:"cost" validate:"gte=0"`
	Quantity int    `json:"quantity" validate:"gte=1"`
}

// UpdateProduct defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateProduct struct {
	Name     *string `json:"name"`
	Cost     *int    `json:"cost" validate:"omitempty,gte=0"`
	Quantity *int    `json:"quantity" validate:"omitempty,gte=1"`
}

// List gets all Products from the database.
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	var products []Product

	// TODO: Talk about issues of using '*' after talking about migrations.
	if err := db.SelectContext(ctx, &products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Create uses the provided *Product to insert a new product record. Generated
// fields like ID, DateCreated, and DateUpdated are populated.
func Create(ctx context.Context, db *sqlx.DB, n NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        n.Name,
		Cost:        n.Cost,
		Quantity:    n.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO products
		(product_id, name, cost, quantity, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated,
	)
	if err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}

	return &p, nil
}

// Get finds the product identified by a given ID.
func Get(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {
	var p Product

	err := db.GetContext(ctx, &p, `
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

// Update modifies fields about a Product. It will error if the specified ID
// does not reference an existing Product.
func Update(ctx context.Context, db *sqlx.DB, id string, update UpdateProduct, now time.Time) error {
	old, err := Get(ctx, db, id)
	if err != nil {
		return err
	}

	if update.Name != nil {
		old.Name = *update.Name
	}
	if update.Cost != nil {
		old.Cost = *update.Cost
	}
	if update.Quantity != nil {
		old.Quantity = *update.Quantity
	}
	old.DateUpdated = now

	_, err = db.ExecContext(ctx, `
		UPDATE products SET
			"name" = $2,
			"cost" = $3,
			"quantity" = $4,
			"date_updated" = $5
		WHERE product_id = $1`,
		id, old.Name, old.Cost, old.Quantity, old.DateUpdated,
	)

	if err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}

// Delete removes the product identified by a given ID.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {

	_, err := db.ExecContext(ctx, `
		DELETE FROM products
		WHERE product_id = $1`,
		id,
	)
	if err != nil {
		return errors.Wrap(err, "deleting product")
	}

	return nil
}
