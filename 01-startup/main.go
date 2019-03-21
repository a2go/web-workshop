package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8000", http.HandlerFunc(ListProducts)); err != nil {
		log.Fatalf("error: listening and serving: %s", err)
	}
}

// Product is an item we sell.
type Product struct {
	Name     string
	Cost     int
	Quantity int
}

// ListProducts is an HTTP Handler for returning a list of Products.
func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{
		{Name: "Comic Books", Cost: 50, Quantity: 42},
		{Name: "McDonalds Toys", Cost: 75, Quantity: 120},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("error: encoding response: %s", err)
	}
}
