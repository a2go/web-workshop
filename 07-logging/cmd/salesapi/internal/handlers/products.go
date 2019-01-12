package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/ardanlabs/service-training/07-logging/internal/platform/log"
	"github.com/ardanlabs/service-training/07-logging/internal/products"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	DB *sqlx.DB
}

// List gets all products from the service layer and encodes them for the
// client response.
func (s *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := products.List(s.DB)
	if err != nil {
		log.Log("listing products", "error", err)
		w.WriteHeader(500)
		return
	}

	// TODO: Don't return an array (return an object with an array).
	//       Make a named response type.
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Log("encoding response", "error", err)
		return
	}
}
