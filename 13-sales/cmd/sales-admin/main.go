package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/ardanlabs/service-training/13-sales/internal/platform/database"
	"github.com/ardanlabs/service-training/13-sales/internal/platform/log"
	"github.com/ardanlabs/service-training/13-sales/internal/platform/database/schema"
)

// This is for parsing the environment.
const envKey = "sales"

type config struct {
	DB database.Config
}

func main() {
	if err := run(); err != nil {
		log.Log("shutting down", "error", err)
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
	db, err := database.Open(cfg.DB)
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db.DB); err != nil {
			return errors.Wrap(err, "applying migrations")
		}
		log.Log("Migrations complete")

	case "seed":
		if err := schema.Seed(db.DB); err != nil {
			return errors.Wrap(err, "seeding database")
		}
		log.Log("Seed data complete")
	}

	return nil
}
