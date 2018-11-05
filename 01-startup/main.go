package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func main() {
	if err := http.ListenAndServe(":8000", http.HandlerFunc(ListProducts)); err != nil {
		log.Fatal(errors.Wrap(err, "listening and serving"))
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
		log.Println(errors.Wrap(err, "encoding response"))
	}
}
