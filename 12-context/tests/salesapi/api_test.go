package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ardanlabs/service-training/12-context/internal/products"
	"github.com/kelseyhightower/envconfig"
)

func TestAPI(t *testing.T) {
	var cfg struct {
		HTTP struct {
			Address string `default:"localhost:8000"`
		} `envconfig:"http"`
	}
	envconfig.MustProcess("test", &cfg)

	c := http.Client{
		Timeout: 15 * time.Second,
	}

	// Create a product.
	var p0created products.Product
	t.Run("create", func(t *testing.T) {
		resp, err := c.Post(fmt.Sprintf("http://%s/v1/products", cfg.HTTP.Address), "application/json", strings.NewReader(`{
			"name": "product0",
			"cost": 55,
			"quantity": 6
		}`))
		if err != nil {
			t.Fatalf("posting: %s", err)
		}
		if exp, got := http.StatusOK, resp.StatusCode; exp != got {
			t.Fatalf("posting: expected status code %v, got %v", exp, got)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&p0created); err != nil {
			t.Fatalf("decoding: %s", err)
		}
		if p0created.ID == "" {
			t.Fatal("expected non-empty product id")
		}
		// TODO: More product field assertions.
	})

	// Get the created product.
	var p0gotten products.Product
	t.Run("get", func(t *testing.T) {
		resp, err := c.Get(fmt.Sprintf("http://%s/v1/products/%s", cfg.HTTP.Address, p0created.ID))
		if err != nil {
			t.Fatalf("getting: %s", err)
		}
		if exp, got := http.StatusOK, resp.StatusCode; exp != got {
			t.Fatalf("getting: expected status code %v, got %v", exp, got)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&p0gotten); err != nil {
			t.Fatalf("decoding: %s", err)
		}
		if exp, got := p0created, p0gotten; exp != got {
			t.Fatalf("expected %+v, got %+v", exp, got)
		}
	})

	// List products.
	t.Run("list", func(t *testing.T) {
		resp, err := c.Get(fmt.Sprintf("http://%s/v1/products", cfg.HTTP.Address))
		if err != nil {
			t.Fatalf("getting: %s", err)
		}
		if exp, got := http.StatusOK, resp.StatusCode; exp != got {
			t.Fatalf("getting: expected status code %v, got %v", exp, got)
		}
		defer resp.Body.Close()

		var list []products.Product
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			t.Fatalf("decoding: %s", err)
		}
		if len(list) == 0 {
			t.Fatal("expected non-empty product list")
		}
	})
}
