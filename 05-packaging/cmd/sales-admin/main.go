package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/garagesale/internal/platform/database"
	"github.com/ardanlabs/garagesale/internal/schema"
)

func main() {

	flag.Parse()

	// Initialize dependencies.
	db, err := database.Open()
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Println("error applying migrations", err)
			os.Exit(1)
		}
		fmt.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Println("error seeding database", err)
			os.Exit(1)
		}
		fmt.Println("Seed data complete")
		return
	}
}
