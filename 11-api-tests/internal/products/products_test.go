package products_test

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/ardanlabs/service-training/11-api-tests/internal/products"
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
		if *p1 != *p1 {
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
		DB struct {
			User     string `default:"postgres"`
			Password string `default:"postgres"`
			Host     string `default:"localhost"`
			Name     string `default:"postgres"`
		}
	}
	envconfig.MustProcess("TEST", &cfg)

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DB.User, cfg.DB.Password),
		Host:   cfg.DB.Host,
		Path:   cfg.DB.Name,
		RawQuery: (url.Values{
			"sslmode":  []string{"disable"},
			"timezone": []string{"utc"},
		}).Encode(),
	}

	db0, err := sqlx.Connect("postgres", u.String())
	if err != nil {
		t.Fatalf("connecting to db: %s", err)
	}

	newDB := fmt.Sprintf("%v_%v", strings.ToLower(t.Name()), time.Now().UnixNano())
	if _, err := db0.Exec("CREATE DATABASE " + newDB); err != nil {
		t.Fatalf("creating database %q: %s", newDB, err)
	}

	u.Path = newDB
	db, err := sqlx.Connect("postgres", u.String())
	if err != nil {
		t.Fatalf("connecting to db: %s", err)
	}

	schema, err := ioutil.ReadFile("../../schema.sql")
	if err != nil {
		t.Fatalf("reading schema file: %s", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("migrating: %s", err)
	}

	return db, func() {
		// Cleanup
		db.Close()
		if _, err := db0.Exec("DROP DATABASE " + newDB); err != nil {
			t.Errorf("dropping database: %s", err)
		}
		db0.Close()
	}
}
