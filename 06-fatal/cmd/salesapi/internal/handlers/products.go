package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ardanlabs/service-training/06-fatal/internal/products"
	"github.com/jmoiron/sqlx"
)

type Products struct {
	DB *sqlx.DB
}

func (s *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := products.List(s.DB)
	if err != nil {
		log.Printf("error: listing products: %s", err)
		w.WriteHeader(500)
		return
	}

	// TODO: Don't return an array (return an object with an array).
	//       Make a named response type.
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Printf("error: encoding response: %s", err)
		return
	}
}
