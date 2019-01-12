package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	// TODO: Mention the idiosyncrasies of using the sql pkg.
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ardanlabs/service-training/03-database/schema"
)

// 1. Start postgres:
// docker-compose up -d
//
// 2. Create the schema and insert some seed data.
// go build
// ./03-database migrate
// ./03-database seed
//
// 3. Run the app then make requests.
// ./03-database

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

	service := Service{db: db}

	server := http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(service.ListProducts),
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

// TODO: Mention JSON conventions / consistency and `json` tags in later (API) session.

// Product is an item we sell.
type Product struct {
	ID       string `db:"product_id"`
	Name     string `db:"name"`
	Cost     int    `db:"cost"`
	Quantity int    `db:"quantity"`
}

// Service holds business logic related to Products.
type Service struct {
	db *sqlx.DB
}

// ListProducts gets all Products from the database then encodes them in a
// response to the client.
func (s *Service) ListProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product

	// TODO: Seperate layers of concern in later section.
	// TODO: Talk about issues of using '*'.
	if err := s.db.Select(&products, "SELECT * FROM products"); err != nil {
		log.Printf("error: selecting products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("error: encoding response: %s", err)
		return
	}
}
