package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/web"
	"github.com/ardanlabs/garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger
}

// List gets all products from the service layer.
func (p *Products) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Product.List")
	defer span.End()

	list, err := product.List(ctx, p.db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (p *Products) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Products.Create")
	defer span.End()

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	prod, err := product.Create(ctx, p.db, claims, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Respond(ctx, w, &prod, http.StatusCreated)
}

// Retrieve finds a single product identified by an ID in the request URL.
func (p *Products) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Products.Get")
	defer span.End()

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(ctx, p.db, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.RespondError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.RespondError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting product %q", id)
		}
	}

	return web.Respond(ctx, w, prod, http.StatusOK)
}

// Update decodes the body of a request to update an existing product. The ID
// of the product is part of the request URL.
func (p *Products) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Products.Update")
	defer span.End()

	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	if err := product.Update(ctx, p.db, claims, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.RespondError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.RespondError(err, http.StatusBadRequest)
		case product.ErrForbidden:
			return web.RespondError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a single product identified by an ID in the request URL.
func (p *Products) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Products.Delete")
	defer span.End()

	id := chi.URLParam(r, "id")

	if err := product.Delete(ctx, p.db, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.RespondError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (p *Products) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Products.AddSale")
	defer span.End()

	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(ctx, p.db, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(ctx, w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (p *Products) ListSales(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Products.ListSales")
	defer span.End()

	id := chi.URLParam(r, "id")

	list, err := product.ListSales(ctx, p.db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}
