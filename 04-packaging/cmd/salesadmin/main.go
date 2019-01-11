package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

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
}
