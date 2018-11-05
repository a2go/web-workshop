package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/service-training/04-packaging/cmd/salesapi/internal/handlers"
	"github.com/ardanlabs/service-training/04-packaging/internal/products"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func main() {
	// Initialize dependencies.
	db, err := sqlx.Connect("postgres", products.DBConn("postgres", "postgres", "localhost", "postgres", true))
	if err != nil {
		log.Fatal(errors.Wrap(err, "connecting to db"))
	}
	defer db.Close()

	productsHandler := handlers.Products{DB: db}
	server := http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(productsHandler.List),
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	log.Print("startup complete")

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
