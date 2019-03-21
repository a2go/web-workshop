package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/ardanlabs/garagesale/internal/products"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger
}

// List gets all products from the service layer and encodes them for the
// client response.
func (s *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := products.List(s.db)
	if err != nil {
		s.log.Println("listing products", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// TODO: Don't return an array (return an object with an array).
	//       Make a named response type.
	if err := json.NewEncoder(w).Encode(list); err != nil {
		s.log.Println("encoding response", "error", err)
		return
	}
}
