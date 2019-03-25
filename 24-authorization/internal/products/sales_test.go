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

	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Create two products to work with.
	newPuzzles := products.NewProduct{
		Name:     "Puzzles",
		Cost:     25,
		Quantity: 6,
	}

	puzzles, err := products.Create(context.Background(), db, newPuzzles, now)
	if err != nil {
		t.Fatalf("creating product: %s", err)
	}

	newToys := products.NewProduct{
		Name:     "Toys",
		Cost:     40,
		Quantity: 3,
	}
	toys, err := products.Create(context.Background(), db, newToys, now)
	if err != nil {
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
