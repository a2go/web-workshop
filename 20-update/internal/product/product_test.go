package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/a2go/garagesale/internal/product"
	"github.com/a2go/garagesale/internal/schema"
	"github.com/a2go/garagesale/internal/tests"
	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	newP := product.NewProduct{
		Name:     "Comic Book",
		Cost:     10,
		Quantity: 55,
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	p0, err := product.Create(ctx, db, newP, now)
	if err != nil {
		t.Fatalf("creating product p0: %s", err)
	}

	p1, err := product.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

	update := product.UpdateProduct{
		Name: tests.StringPointer("Comics"),
		Cost: tests.IntPointer(25),
	}
	updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)

	if err := product.Update(ctx, db, p0.ID, update, updatedTime); err != nil {
		t.Fatalf("creating product p0: %s", err)
	}

	saved, err := product.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}

	// Check specified fields were updated. Make a copy of the original product
	// and change just the fields we expect then diff it with what was saved.
	want := *p0
	want.Name = "Comics"
	want.Cost = 25
	want.DateUpdated = updatedTime

	if diff := cmp.Diff(want, *saved); diff != "" {
		t.Fatalf("updated record did not match:\n%s", diff)
	}
}

func TestProductList(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := product.List(context.Background(), db)
	if err != nil {
		t.Fatalf("listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
