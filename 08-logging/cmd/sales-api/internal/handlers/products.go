package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ardanlabs/garagesale/internal/products"
	"github.com/jmoiron/sqlx"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List gets all products from the service layer and encodes them for the
// client response.
func (s *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := products.List(s.DB)
	if err != nil {
		s.Log.Printf("error: listing products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		s.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		s.Log.Println("error writing result", err)
	}
}
