package handlers

import (
	"net/http"

	"github.com/ardanlabs/service-training/09-refactor-handlers/internal/platform/web"
	"github.com/ardanlabs/service-training/09-refactor-handlers/internal/products"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Products struct {
	db *sqlx.DB

	http.Handler
}

func NewProducts(db *sqlx.DB) *Products {
	p := Products{db: db}

	r := chi.NewRouter()
	r.Post("/v1/products", web.Encode(p.Create))
	r.Get("/v1/products", web.Encode(p.List))
	r.Get("/v1/products/{id}", web.Encode(p.Get))
	p.Handler = r

	return &p
}

func (s *Products) Create(r *http.Request) (interface{}, error) {
	var p products.Product
	if err := web.Decode(r, &p); err != nil {
		return nil, err
	}

	if err := products.Create(s.db, &p); err != nil {
		return nil, errors.Wrap(err, "creating")
	}

	return p, nil
}

func (s *Products) List(r *http.Request) (interface{}, error) {
	list, err := products.List(s.db)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Products) Get(r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")

	p, err := products.Get(s.db, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}
