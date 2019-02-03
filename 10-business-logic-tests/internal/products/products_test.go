package products_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"

	"github.com/ardanlabs/service-training/10-business-logic-tests/internal/platform/database"
	"github.com/ardanlabs/service-training/10-business-logic-tests/internal/products"
	"github.com/ardanlabs/service-training/10-business-logic-tests/internal/schema"
)

func TestProducts(t *testing.T) {
	db, drop := testDB(t)
	defer drop()

	{ // Create and Get.
		p0 := products.Product{
			Name:     "Comic Book",
			Cost:     10,
			Quantity: 55,
		}
		if err := products.Create(db, &p0); err != nil {
			t.Fatalf("creating product p0: %s", err)
		}
		p1, err := products.Get(db, p0.ID)
		if err != nil {
			t.Fatalf("getting product p0: %s", err)
		}
		if *p1 != p0 {
			t.Fatalf("fetched != created: %v != %v", p1, p0)
		}
	}

	{ // List.
		ps, err := products.List(db)
		if err != nil {
			t.Fatalf("listing products: %s", err)
		}
		if exp, got := 1, len(ps); exp != got {
			t.Fatalf("expected product list size %v, got %v", exp, got)
		}
	}
}

func testDB(t *testing.T) (*sqlx.DB, func()) {
	var cfg struct {
		DB database.Config
	}
	envconfig.MustProcess("TEST", &cfg)

	db0, err := database.Open(cfg.DB)
	if err != nil {
		t.Fatalf("connecting to db: %s", err)
	}

	newDB := fmt.Sprintf("%v_test_%v", strings.ToLower(t.Name()), time.Now().UnixNano())
	if _, err := db0.Exec("CREATE DATABASE " + newDB); err != nil {
		t.Fatalf("creating database %q: %s", newDB, err)
	}

	cfg.DB.Name = newDB
	db, err := database.Open(cfg.DB)
	if err != nil {
		t.Fatalf("connecting to db: %s", err)
	}

	if err := schema.Migrate(db.DB); err != nil {
		t.Fatalf("migrating: %s", err)
	}

	// cleanup is the function that should be invoked when the caller is done
	// with the database.
	cleanup := func() {
		db.Close()
		if _, err := db0.Exec("DROP DATABASE " + newDB); err != nil {
			t.Errorf("dropping database: %s", err)
		}
		db0.Close()
	}

	return db, cleanup
}
