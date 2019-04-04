package products_test

import (
	"testing"
	"time"

	"github.com/ardanlabs/garagesale/internal/platform/database/databasetest"
	"github.com/ardanlabs/garagesale/internal/platform/database/schema"
	"github.com/ardanlabs/garagesale/internal/products"
	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	newP := products.NewProduct{
		Name:     "Comic Book",
		Cost:     10,
		Quantity: 55,
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	p0, err := products.Create(db, newP, now)
	if err != nil {
		t.Fatalf("creating product p0: %s", err)
	}

	p1, err := products.Get(db, p0.ID)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}
}

func TestProductList(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	if err := schema.Seed(db.DB); err != nil {
		t.Fatal(err)
	}

	ps, err := products.List(db)
	if err != nil {
		t.Fatalf("listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
