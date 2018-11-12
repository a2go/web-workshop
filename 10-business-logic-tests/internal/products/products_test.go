package products_test

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ardanlabs/service-training/10-business-logic-tests/internal/products"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func TestProducts(t *testing.T) {
	db, drop := testDB(t)
	defer drop()

	{ // Create and Get.
		p0 := products.Product{}
		if err := products.Create(db, &p0); err != nil {
			t.Fatal(errors.Wrap(err, "creating product p0"))
		}
		p1, err := products.Get(db, p0.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "getting product p0"))
		}
		if *p1 != *p1 {
			t.Fatalf("fetched != created: %v != %v", p1, p0)
		}
	}

	{ // List.
		ps, err := products.List(db)
		if err != nil {
			t.Fatal(errors.Wrap(err, "listing products"))
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
		t.Fatal(errors.Wrap(err, "connecting to db"))
	}

	newDB := fmt.Sprintf("%v_%v", strings.ToLower(t.Name()), time.Now().UnixNano())
	db0.Exec("CREATE DATABASE " + newDB)

	u.Path = newDB
	db, err := sqlx.Connect("postgres", u.String())
	if err != nil {
		t.Fatal(errors.Wrap(err, "connecting to db"))
	}

	schema, err := ioutil.ReadFile("../../schema.sql")
	if err != nil {
		t.Fatal(errors.Wrap(err, "reading schema file"))
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatal(errors.Wrap(err, "migrating"))
	}

	return db, func() {
		// Cleanup
		db.Close()
		db0.Exec("DROP DATABASE " + newDB)
		db0.Close()
	}
}
