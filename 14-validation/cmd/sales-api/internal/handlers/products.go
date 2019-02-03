package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ardanlabs/service-training/14-validation/internal/platform/web"
	"github.com/ardanlabs/service-training/14-validation/internal/products"
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
	r.Post("/v1/products", web.Encode(p.Create))
	r.Get("/v1/products", web.Encode(p.List))
	r.Get("/v1/products/{id}", web.Encode(p.Get))

	r.Post("/v1/products/{id}/sales", web.Encode(p.AddSale))
	r.Get("/v1/products/{id}/sales", web.Encode(p.ListSales))
	p.Handler = r

	return &p
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (s *Products) Create(r *http.Request) (interface{}, error) {
	var p products.Product
	if err := web.Decode(r, &p); err != nil {
		return nil, err
	}

	if err := products.Create(r.Context(), s.db, &p); err != nil {
		return nil, errors.Wrap(err, "creating")
	}

	return p, nil
}

// List gets all products from the service layer.
func (s *Products) List(r *http.Request) (interface{}, error) {
	list, err := products.List(r.Context(), s.db)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Get finds a single product identified by an ID in the request URL.
func (s *Products) Get(r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")

	p, err := products.Get(r.Context(), s.db, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (s *Products) AddSale(r *http.Request) (interface{}, error) {
	var sale products.Sale
	if err := web.Decode(r, &s); err != nil {
		return nil, err
	}

	sale.ProductID = chi.URLParam(r, "id")

	if err := products.AddSale(r.Context(), s.db, &sale); err != nil {
		return nil, errors.Wrap(err, "creating")
	}

	return sale, nil
}

// ListSales gets all sales for a particular product.
func (s *Products) ListSales(r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")

	list, err := products.ListSales(r.Context(), s.db, id)
	if err != nil {
		return nil, err
	}

	return list, nil
}
