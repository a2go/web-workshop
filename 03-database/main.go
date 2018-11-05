package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	// TODO: Mention the idiosyncrasies of using the sql pkg.
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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
	user, pass, host, name := "postgres", "postgres", "localhost", "postgres"
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&timezone=utc",
		user, pass, host, name))
	if err != nil {
		log.Fatal(errors.Wrap(err, "connecting to db"))
	}
	defer db.Close()

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
		log.Println(errors.Wrap(err, "selecting products"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Println(errors.Wrap(err, "encoding response"))
		return
	}
}
