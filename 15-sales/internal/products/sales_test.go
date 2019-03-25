package products_test

import (
	"context"
	"testing"

	"github.com/ardanlabs/garagesale/internal/platform/database/databasetest"
	"github.com/ardanlabs/garagesale/internal/products"
)

func TestSales(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	// Create two products to work with.
	puzzles := products.Product{
		Name:     "Puzzles",
		Cost:     25,
		Quantity: 6,
	}
	if err := products.Create(context.Background(), db, &puzzles); err != nil {
		t.Fatalf("creating product: %s", err)
	}
	toys := products.Product{
		Name:     "Toys",
		Cost:     40,
		Quantity: 3,
	}
	if err := products.Create(context.Background(), db, &toys); err != nil {
		t.Fatalf("creating product: %s", err)
	}

	{ // Add and list

		s := products.Sale{
			ProductID: puzzles.ID,
			Quantity:  3,
			Paid:      70,
		}

		if err := products.AddSale(context.Background(), db, &s); err != nil {
			t.Fatalf("adding sale: %s", err)
		}

		// Puzzles should show the 1 sale.
		sales, err := products.ListSales(context.Background(), db, puzzles.ID)
		if err != nil {
			t.Fatalf("listing sales: %s", err)
		}
		if exp, got := 1, len(sales); exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}

		if exp, got := s.ID, sales[0].ID; exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}

		// Toys should have 0 sales.
		sales, err = products.ListSales(context.Background(), db, toys.ID)
		if err != nil {
			t.Fatalf("listing sales: %s", err)
		}
		if exp, got := 0, len(sales); exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}
	}
}
