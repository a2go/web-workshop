package products_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/ardanlabs/service-training/16-crud/internal/platform/database/databasetest"
	"github.com/ardanlabs/service-training/16-crud/internal/platform/database/schema"
	"github.com/ardanlabs/service-training/16-crud/internal/platform/tests"
	"github.com/ardanlabs/service-training/16-crud/internal/products"
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
	p0, err := products.Create(context.Background(), db, newP, now)
	if err != nil {
		t.Fatalf("creating product p0: %s", err)
	}

	p1, err := products.Get(context.Background(), db, p0.ID)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

	update := products.UpdateProduct{
		Name: tests.StringPointer("Comics"),
		Cost: tests.IntPointer(25),
	}
	updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)

	if err := products.Update(context.Background(), db, p0.ID, update, updatedTime); err != nil {
		t.Fatalf("creating product p0: %s", err)
	}

	saved, err := products.Get(context.Background(), db, p0.ID)
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

	if err := products.Delete(context.Background(), db, p0.ID); err != nil {
		t.Fatalf("deleting product: %v", err)
	}

	_, err = products.Get(context.Background(), db, p0.ID)
	if err == nil {
		t.Fatalf("should not be able to retrieve deleted product")
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
