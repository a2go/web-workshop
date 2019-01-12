package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"

	"github.com/ardanlabs/service-training/08-crud/internal/platform/log"
	"github.com/ardanlabs/service-training/08-crud/internal/products"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db *sqlx.DB

	http.Handler
}

// NewProducts creates a product handler with multiple routes defined.
func NewProducts(db *sqlx.DB) *Products {
	p := Products{db: db}

	r := chi.NewRouter()
	r.Post("/v1/products", p.Create)
	r.Get("/v1/products", p.List)
	r.Get("/v1/products/{id}", p.Get)
	p.Handler = r

	return &p
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (s *Products) Create(w http.ResponseWriter, r *http.Request) {
	var p products.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Log("decoding product", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := products.Create(s.db, &p); err != nil {
		log.Log("creating product", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Log("encoding response", "error", err)
		return
	}
}

// List gets all products from the service layer and encodes them for the
// client response.
func (s *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := products.List(s.db)
	if err != nil {
		log.Log("listing products", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Don't return an array (return an object with an array).
	//       Make a named response type.
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Log("encoding response", "error", err)
		return
	}
}

// Get finds a single product identified by an ID in the request URL.
func (s *Products) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := products.Get(s.db, id)
	if err != nil {
		log.Log("getting product", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Log("encoding response", "error", err)
		return
	}
}
