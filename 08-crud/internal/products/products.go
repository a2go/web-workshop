package products

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Product struct {
	ID       string `db:"product_id"`
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
	// TODO: Catch "null" in tests with empty list.
	var products []Product

	// TODO: Talk about issues of using '*' after talking about migrations.
	if err := db.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

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

func Get(db *sqlx.DB, id string) (*Product, error) {
	var p Product

	err := db.Get(&p, `
		SELECT * FROM products
		WHERE product_id = $1
		LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, errors.Wrap(err, "selecting single product")
	}

	return &p, nil
}
