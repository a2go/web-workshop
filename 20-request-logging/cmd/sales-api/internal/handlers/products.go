package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/ardanlabs/garagesale/internal/products"
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
	var newP products.NewProduct
	if err := web.Decode(r, &newP); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	p, err := products.Create(r.Context(), s.db, newP, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Encode(r.Context(), w, &p, http.StatusCreated)
}

// List gets all products from the service layer.
func (s *Products) List(w http.ResponseWriter, r *http.Request) error {
	list, err := products.List(r.Context(), s.db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Encode(r.Context(), w, list, http.StatusOK)
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

	return web.Encode(r.Context(), w, p, http.StatusOK)
}

// Update decodes the body of a request to update an existing product a new
// product. The ID of the product is part of the request URL.
func (s *Products) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update products.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := products.Update(r.Context(), s.db, id, update, time.Now()); err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Encode(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single product identified by an ID in the request URL.
func (s *Products) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := products.Delete(r.Context(), s.db, id); err != nil {
		return errors.Wrapf(err, "getting product %q", id)
	}

	return web.Encode(r.Context(), w, nil, http.StatusNoContent)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (s *Products) AddSale(w http.ResponseWriter, r *http.Request) error {
	var sale products.Sale
	if err := web.Decode(r, &s); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	sale.ProductID = chi.URLParam(r, "id")

	if err := products.AddSale(r.Context(), s.db, &sale); err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Encode(r.Context(), w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (s *Products) ListSales(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := products.ListSales(r.Context(), s.db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Encode(r.Context(), w, list, http.StatusOK)
}
