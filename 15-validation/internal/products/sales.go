package products

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Sale represents one item of a transaction where some amount of a product was
// sold. Quantity is the number of units sold and Paid is the total price paid.
// Note that due to haggling the Paid value might not equal Quantity sold *
// Product cost.
type Sale struct {
	ID        string `db:"sale_id" json:"id"`
	ProductID string `db:"product_id" json:"product_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
	Paid      int    `db:"paid" json:"paid"`
}

// AddSale records a sales transaction for a single Product.
func AddSale(ctx context.Context, db *sqlx.DB, s *Sale) error {
	s.ID = uuid.New().String()

	_, err := db.ExecContext(ctx, `
		INSERT INTO sales
		(sale_id, product_id, quantity, paid)
		VALUES ($1, $2, $3, $4)`,
		s.ID, s.ProductID, s.Quantity, s.Paid,
	)
	if err != nil {
		return errors.Wrap(err, "inserting sale")
	}

	return nil
}

// ListSales gives all Sales for a Product.
func ListSales(ctx context.Context, db *sqlx.DB, productID string) ([]Sale, error) {
	var sales []Sale

	if err := db.SelectContext(ctx, &sales, "SELECT * FROM sales"); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}

	return sales, nil
}
