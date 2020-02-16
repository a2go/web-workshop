package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	// NOTE: Models should not be imported, we want to test the exact JSON. We
	// make the comparison process easier using the go-cmp library.
	"github.com/ardanlabs/garagesale/cmd/sales-api/internal/handlers"
	"github.com/ardanlabs/garagesale/internal/tests"
	"github.com/google/go-cmp/cmp"
)

// TestProducts runs a series of tests to exercise Product behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestProducts(t *testing.T) {
	test := tests.New(t)
	defer test.Teardown()

	shutdown := make(chan os.Signal, 1)
	tests := ProductTests{
		app:        handlers.API(shutdown, test.DB, test.Log, test.Authenticator),
		adminToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("List", tests.List)
	t.Run("CreateRequiresFields", tests.CreateRequiresFields)
	t.Run("ProductCRUD", tests.ProductCRUD)
}

// ProductTests holds methods for each product subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type ProductTests struct {
	app        http.Handler
	adminToken string
}

func (p *ProductTests) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/products", nil)
	resp := httptest.NewRecorder()

	req.Header.Set("Authorization", "Bearer "+p.adminToken)

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := []map[string]interface{}{
		{
			"id":           "a2b0639f-2cc6-44b8-b97b-15d69dbb511e",
			"name":         "Comic Books",
			"cost":         float64(50),
			"quantity":     float64(42),
			"revenue":      float64(350),
			"sold":         float64(7),
			"user_id":      "00000000-0000-0000-0000-000000000000",
			"date_created": "2019-01-01T00:00:01.000001Z",
			"date_updated": "2019-01-01T00:00:01.000001Z",
		},
		{
			"id":           "72f8b983-3eb4-48db-9ed0-e45cc6bd716b",
			"name":         "McDonalds Toys",
			"cost":         float64(75),
			"quantity":     float64(120),
			"revenue":      float64(225),
			"sold":         float64(3),
			"user_id":      "00000000-0000-0000-0000-000000000000",
			"date_created": "2019-01-01T00:00:02.000001Z",
			"date_updated": "2019-01-01T00:00:02.000001Z",
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

func (p *ProductTests) CreateRequiresFields(t *testing.T) {
	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("POST", "/v1/products", body)

	req.Header.Set("Authorization", "Bearer "+p.adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusBadRequest, resp.Code)
	}
}

func (p *ProductTests) ProductCRUD(t *testing.T) {
	var created map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"name":"product0","cost":55,"quantity":6}`)

		req := httptest.NewRequest("POST", "/v1/products", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusCreated != resp.Code {
			t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatal("expected non-empty product id")
		}
		if created["date_created"] == "" || created["date_created"] == nil {
			t.Fatal("expected non-empty product date_created")
		}
		if created["date_updated"] == "" || created["date_updated"] == nil {
			t.Fatal("expected non-empty product date_updated")
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": created["date_updated"],
			"name":         "product0",
			"cost":         float64(55),
			"quantity":     float64(6),
			"sold":         float64(0),
			"revenue":      float64(0),
			"user_id":      tests.AdminID,
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("Response did not match expected. Diff:\n%s", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		// Fetched product should match the one we created.
		if diff := cmp.Diff(created, fetched); diff != "" {
			t.Fatalf("Retrieved product should match created. Diff:\n%s", diff)
		}
	}

	{ // UPDATE
		body := strings.NewReader(`{"name":"new name","cost":20,"quantity":10}`)
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("PUT", url, body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusNoContent != resp.Code {
			t.Fatalf("updating: expected status code %v, got %v", http.StatusNoContent, resp.Code)
		}

		// Retrieve updated record to be sure it worked.
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp = httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var updated map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": updated["date_updated"],
			"name":         "new name",
			"cost":         float64(20),
			"quantity":     float64(10),
			"sold":         float64(0),
			"revenue":      float64(0),
			"user_id":      tests.AdminID,
		}

		// Updated product should match the one we created.
		if diff := cmp.Diff(want, updated); diff != "" {
			t.Fatalf("Retrieved product should match created. Diff:\n%s", diff)
		}
	}

	{ // DELETE
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusNoContent != resp.Code {
			t.Fatalf("updating: expected status code %v, got %v", http.StatusNoContent, resp.Code)
		}

		// Retrieve updated record to be sure it worked.
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp = httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusNotFound, resp.Code)
		}
	}
}
