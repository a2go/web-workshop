package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/garagesale/internal/platform/database"
	"github.com/ardanlabs/garagesale/internal/schema"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// This is for parsing the environment.
const envKey = "sales"

type config struct {
	DB struct {
		User       string `default:"postgres"`
		Password   string `default:"postgres" json:"-"` // Prevent the marshalling of secrets.
		Host       string `default:"localhost"`
		Name       string `default:"postgres"`
		DisableTLS bool   `default:"false" split_words:"true"`
	}
}

func main() {
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {

	// Process command line flags.
	var flags struct {
		configOnly bool
	}
	flag.BoolVar(&flags.configOnly, "config-only", false, "only show parsed configuration and exit")
	flag.Usage = func() {
		fmt.Print("This program is a CLI tool for administering the Garage Sale service.\n\nUsage of sales-admin:\n\nsales-admin [flags]\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
		fmt.Print("\nConfiguration:\n\n")
		envconfig.Usage(envKey, &config{})
	}
	flag.Parse()

	// Get configuration from environment.
	var cfg config
	if err := envconfig.Process(envKey, &cfg); err != nil {
		return errors.Wrap(err, "parsing config")
	}

	// Print config and exit if requested.
	if flags.configOnly {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "	")
		if err := enc.Encode(cfg); err != nil {
			return errors.Wrap(err, "encoding config as json")
		}
		return nil
	}

	// Initialize dependencies.
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "applying migrations")
		}
		fmt.Println("Migrations complete")

	case "seed":
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "seeding database")
		}
		fmt.Println("Seed data complete")
	}

	return nil
}
