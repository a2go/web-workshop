package products_test

import (
	"context"
	"testing"
	"time"

	"github.com/ardanlabs/garagesale/internal/platform/database/databasetest"
	"github.com/ardanlabs/garagesale/internal/products"
)

func TestSales(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	// Create a product for Sales to work with.
	newP := products.NewProduct{
		Name:     "Puzzles",
		Cost:     25,
		Quantity: 6,
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	p, err := products.Create(context.Background(), db, newP, now)
	if err != nil {
		t.Fatalf("creating product: %s", err)
	}

	{ // Add and list

		s := products.Sale{
			ProductID: p.ID,
			Quantity:  3,
			Paid:      75,
		}

		if err := products.AddSale(context.Background(), db, &s); err != nil {
			t.Fatalf("adding sale: %s", err)
		}

		sales, err := products.ListSales(context.Background(), db, p.ID)
		if err != nil {
			t.Fatalf("listing sales: %s", err)
		}
		if exp, got := 1, len(sales); exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}

		if exp, got := s.ID, sales[0].ID; exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}
	}
}
