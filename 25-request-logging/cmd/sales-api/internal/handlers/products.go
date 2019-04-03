package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/ardanlabs/garagesale/internal/products"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger
}

// List gets all products from the service layer.
func (s *Products) List(w http.ResponseWriter, r *http.Request) error {
	list, err := products.List(r.Context(), s.db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (s *Products) Create(w http.ResponseWriter, r *http.Request) error {
	var np products.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	p, err := products.Create(r.Context(), s.db, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Respond(r.Context(), w, &p, http.StatusCreated)
}

// Get finds a single product identified by an ID in the request URL.
func (s *Products) Get(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	p, err := products.Get(r.Context(), s.db, id)
	if err != nil {
		switch err {
		case products.ErrNotFound:
			return web.ErrorWithStatus(err, http.StatusNotFound)
		case products.ErrInvalidID:
			return web.ErrorWithStatus(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting product %q", id)
		}
	}

	return web.Respond(r.Context(), w, p, http.StatusOK)
}

// Update decodes the body of a request to update an existing product. The ID
// of the product is part of the request URL.
func (s *Products) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update products.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := products.Update(r.Context(), s.db, id, update, time.Now()); err != nil {
		switch err {
		case products.ErrNotFound:
			return web.ErrorWithStatus(err, http.StatusNotFound)
		case products.ErrInvalidID:
			return web.ErrorWithStatus(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single product identified by an ID in the request URL.
func (s *Products) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := products.Delete(r.Context(), s.db, id); err != nil {
		switch err {
		case products.ErrInvalidID:
			return web.ErrorWithStatus(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (s *Products) AddSale(w http.ResponseWriter, r *http.Request) error {
	var ns products.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := products.AddSale(r.Context(), s.db, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(r.Context(), w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (s *Products) ListSales(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := products.ListSales(r.Context(), s.db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}
