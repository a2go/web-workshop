package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	// TODO: Mention the idiosyncrasies of using the sql pkg.
	_ "github.com/lib/pq"
)

// 1. Start postgres in another terminal:
// docker run --rm -p 5432:5432 postgres
//
// 2. Connect to postgres:
// docker run --rm -it --network="host" postgres psql -U postgres -h localhost
// > CREATE TABLE products (id SERIAL PRIMARY KEY, name VARCHAR(255), cost INT, quantity INT);
// > INSERT INTO products (name, quantity, cost) VALUES ('Comic Books', 50, 42);
// > INSERT INTO products (name, quantity, cost) VALUES ('McDonalds Toys', 75, 120);
// > \q

func main() {
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
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

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
type Product struct {
	ID       string `db:"product_id"`
	Name     string `db:"name"`
	Cost     int    `db:"cost"`
	Quantity int    `db:"quantity"`
}

type Service struct {
	db *sqlx.DB
}

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
