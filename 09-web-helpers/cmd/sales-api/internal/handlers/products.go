package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"

	"github.com/ardanlabs/service-training/09-web-helpers/internal/platform/web"
	"github.com/ardanlabs/service-training/09-web-helpers/internal/products"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (s *Products) Create(w http.ResponseWriter, r *http.Request) {
	var p products.Product
	if err := web.Decode(r, &p); err != nil {
		s.log.Println("decoding product", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := products.Create(s.db, &p); err != nil {
		s.log.Println("creating product", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Encode(w, &p, http.StatusCreated); err != nil {
		s.log.Println("encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// List gets all products from the service layer.
func (s *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := products.List(s.db)
	if err != nil {
		s.log.Println("listing products", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Encode(w, list, http.StatusOK); err != nil {
		s.log.Println("encoding response", "error", err)
		return
	}
}

// Get finds a single product identified by an ID in the request URL.
func (s *Products) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := products.Get(s.db, id)
	if err != nil {
		s.log.Println("getting product", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Encode(w, p, http.StatusOK); err != nil {
		s.log.Println("encoding response", "error", err)
		return
	}
}
