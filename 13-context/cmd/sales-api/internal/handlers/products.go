package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ardanlabs/service-training/13-context/internal/platform/web"
	"github.com/ardanlabs/service-training/13-context/internal/products"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (s *Products) Create(w http.ResponseWriter, r *http.Request) error {
	var p products.Product
	if err := web.Decode(r, &p); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	if err := products.Create(r.Context(), s.db, &p); err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Encode(w, &p, http.StatusCreated)
}

// List gets all products from the service layer.
func (s *Products) List(w http.ResponseWriter, r *http.Request) error {
	list, err := products.List(r.Context(), s.db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Encode(w, list, http.StatusOK)
}

// Get finds a single product identified by an ID in the request URL.
func (s *Products) Get(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	p, err := products.Get(r.Context(), s.db, id)
	if err != nil {
		if err == products.ErrNotFound {
			return web.ErrorWithStatus(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "getting product %q", id)
	}

	return web.Encode(w, p, http.StatusOK)
}
