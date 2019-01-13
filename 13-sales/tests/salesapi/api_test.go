package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	/* NOTE: Models should not be imported, we want to test the exact JSON */

	"github.com/kelseyhightower/envconfig"
)

func TestAPI(t *testing.T) {
	var cfg struct {
		HTTP struct {
			Address string `default:"localhost:8000"`
		} `envconfig:"http"`
	}
	envconfig.MustProcess("test_api", &cfg)

	c := http.Client{
		Timeout: 15 * time.Second,
	}

	var list0 []map[string]interface{}
	t.Run("GetInitialProductList", func(t *testing.T) {
		resp, err := c.Get(fmt.Sprintf("http://%s/v1/products", cfg.HTTP.Address))
		if err != nil {
			t.Fatalf("getting: %s", err)
		}
		if exp, got := http.StatusOK, resp.StatusCode; exp != got {
			t.Fatalf("getting: expected status code %v, got %v", exp, got)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&list0); err != nil {
			t.Fatalf("decoding: %s", err)
		}
	})

	var p0created map[string]interface{}
	t.Run("CreateNewProduct", func(t *testing.T) {
		tocreate := map[string]interface{}{
			"name":     "product0",
			"cost":     55,
			"quantity": 6,
		}
		resp, err := c.Post(fmt.Sprintf("http://%s/v1/products", cfg.HTTP.Address), "application/json", jsonRdr(t, tocreate))
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

		if p0created["id"] == "" {
			t.Fatal("expected non-empty product id")
		}

		// The created product should be returned in the response.
		assertJSONEquals(t, tocreate, mapWithout(p0created, "id"))
	})

	var p0gotten map[string]interface{}
	t.Run("GetCreated", func(t *testing.T) {
		resp, err := c.Get(fmt.Sprintf("http://%s/v1/products/%s", cfg.HTTP.Address, p0created["id"]))
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

		// Fetched product should match the one we created.
		assertJSONEquals(t, p0created, p0gotten)
	})

	t.Run("ListAfterCreating", func(t *testing.T) {
		resp, err := c.Get(fmt.Sprintf("http://%s/v1/products", cfg.HTTP.Address))
		if err != nil {
			t.Fatalf("getting: %s", err)
		}
		if exp, got := http.StatusOK, resp.StatusCode; exp != got {
			t.Fatalf("getting: expected status code %v, got %v", exp, got)
		}
		defer resp.Body.Close()

		var list []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		// List should be 1 product larger than before.
		if exp, got := len(list0)+1, len(list); exp != got {
			t.Fatalf("expected list length = %v (list0 + 1), got %v", exp, got)
		}

		// List should contain the created product.
		var foundP0created bool
		for _, p := range list {
			if p["id"] == p0created["id"] {
				foundP0created = true
			}
		}
		if !foundP0created {
			t.Fatalf("did not find id %q in list", p0created["id"])
		}
	})
}

// mapWithout returns a new map with a given key.
func mapWithout(m map[string]interface{}, key string) map[string]interface{} {
	newM := make(map[string]interface{})
	for k, v := range m {
		if k != key {
			newM[k] = v
		}
	}
	return newM
}

// assertJSONEquals compares two maps by marshalling them to JSON.
func assertJSONEquals(t *testing.T, exp, got map[string]interface{}) {
	if exp, got := string(toJSON(t, exp)), string(toJSON(t, got)); exp != got {
		t.Fatalf("expected json object:\n%s\nreceived object:\n%s\n", exp, got)
	}
}

// jsonRdr marshals and returns a reader with the resulting JSON.
func jsonRdr(t *testing.T, x interface{}) io.Reader {
	return bytes.NewReader(toJSON(t, x))
}

func toJSON(t *testing.T, x interface{}) []byte {
	btys, err := json.Marshal(x)
	if err != nil {
		t.Fatalf("marshalling to json: %s", err)
	}
	return btys
}
