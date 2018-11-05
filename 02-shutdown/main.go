package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(ListProducts),
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	log.Print("initialized")

	select {
	case err := <-serverErrors:
		log.Fatal(errors.Wrap(err, "listening and serving"))
	case <-osSignals:
		log.Print("caught signal, shutting down")

		// Give outstanding requests 30 seconds to complete.
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error: %s", errors.Wrap(err, "shutting down server"))
			if err := server.Close(); err != nil {
				log.Printf("error: %s", errors.Wrap(err, "forcing server to close"))
			}
		}
	}

	log.Print("done")
}

type Product struct {
	Name     string
	Cost     int
	Quantity int
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	// Simulate a long-running request.
	time.Sleep(10 * time.Second)

	products := []Product{
		{Name: "Comic Books", Cost: 50, Quantity: 42},
		{Name: "McDonalds Toys", Cost: 75, Quantity: 120},
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Println(errors.Wrap(err, "encoding response"))
	}
}
