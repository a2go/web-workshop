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

type Product struct {
	Name     string
	Cost     int
	Quantity int
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{
		{Name: "Comic Books", Cost: 50, Quantity: 42},
		{Name: "McDonalds Toys", Cost: 75, Quantity: 120},
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("error: encoding response: %s", err)
	}
}
