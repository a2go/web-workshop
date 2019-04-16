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
	list, err := products.List(s.db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Respond(w, list, http.StatusOK)
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (s *Products) Create(w http.ResponseWriter, r *http.Request) error {
	var np products.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	p, err := products.Create(s.db, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Respond(w, &p, http.StatusCreated)
}

// Get finds a single product identified by an ID in the request URL.
func (s *Products) Get(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	p, err := products.Get(s.db, id)
	if err != nil {
		switch err {
		case products.ErrNotFound:
			return web.WrapErrorWithStatus(err, http.StatusNotFound)
		case products.ErrInvalidID:
			return web.WrapErrorWithStatus(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting product %q", id)
		}
	}

	return web.Respond(w, p, http.StatusOK)
}
