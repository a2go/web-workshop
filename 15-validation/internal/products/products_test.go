package products_test

import (
	"context"
	"testing"

	"github.com/ardanlabs/garagesale/internal/platform/database/databasetest"
	"github.com/ardanlabs/garagesale/internal/platform/database/schema"
	"github.com/ardanlabs/garagesale/internal/products"
)

func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	p0 := products.Product{
		Name:     "Comic Book",
		Cost:     10,
		Quantity: 55,
	}
	if err := products.Create(context.Background(), db, &p0); err != nil {
		t.Fatalf("creating product p0: %s", err)
	}
	p1, err := products.Get(context.Background(), db, p0.ID)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}
	if *p1 != p0 {
		t.Fatalf("fetched != created: %v != %v", p1, p0)
	}
}

func TestProductList(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	if err := schema.Seed(db.DB); err != nil {
		t.Fatal(err)
	}

	ps, err := products.List(context.Background(), db)
	if err != nil {
		t.Fatalf("listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
