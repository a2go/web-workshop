package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ardanlabs/service-training/04-packaging/cmd/salesapi/internal/handlers"
	"github.com/ardanlabs/service-training/04-packaging/internal/schema"
)

func main() {

	flag.Parse()

	// Initialize dependencies.
	var db *sqlx.DB
	{
		u := url.URL{
			Scheme: "postgres",
			User:   url.UserPassword("postgres", "postgres"),
			Host:   "localhost",
			Path:   "postgres",
			RawQuery: (url.Values{
				"sslmode":  []string{"disable"},
				"timezone": []string{"utc"},
			}).Encode(),
		}

		var err error
		db, err = sqlx.Connect("postgres", u.String())
		if err != nil {
			log.Fatalf("error: connecting to db: %s", err)
		}

		defer db.Close()
	}

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db.DB); err != nil {
			log.Println("error applying migrations", err)
			os.Exit(1)
		}
		log.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db.DB); err != nil {
			log.Println("error seeding database", err)
			os.Exit(1)
		}
		log.Println("Seed data complete")
		return
	}

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
		log.Fatalf("error: listening and serving: %s", err)

	case <-osSignals:
		log.Print("caught signal, shutting down")

		// Give outstanding requests 15 seconds to complete.
		const timeout = 15 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error: gracefully shutting down server: %s", err)
			if err := server.Close(); err != nil {
				log.Printf("error: closing server: %s", err)
			}
		}
	}

	log.Print("done")
}
