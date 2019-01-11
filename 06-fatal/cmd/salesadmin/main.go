package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ardanlabs/service-training/06-fatal/internal/schema"
)

// This is the application name.
const name = "salesapi"

type config struct {
	// NOTE: We don't pass in a connection string b/c our application may assume
	//       certain parameters are set.
	DB struct {
		User     string `default:"postgres"`
		Password string `default:"postgres" json:"-"` // Prevent the marshalling of secrets.
		Host     string `default:"localhost"`
		Name     string `default:"postgres"`

		DisableTLS bool `default:"false" envconfig:"disable_tls"`
	}
}

func main() {
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {
	// Process inputs.
	var flags struct {
		configOnly bool
	}
	flag.Usage = func() {
		fmt.Print("This daemon is a service which manages products.\n\nUsage of sales-api:\n\nsales-api [flags]\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
		fmt.Print("\nConfiguration:\n\n")
		envconfig.Usage(name, &config{})
	}
	flag.BoolVar(&flags.configOnly, "config-only", false, "only show parsed configuration and exit")
	flag.Parse()

	var cfg config
	if err := envconfig.Process(name, &cfg); err != nil {
		return errors.Wrap(err, "parsing config")
	}

	if flags.configOnly {
		if err := json.NewEncoder(os.Stdout).Encode(cfg); err != nil {
			return errors.Wrap(err, "encoding config as json")
		}
		return nil
	}

	// Initialize dependencies.
	var db *sqlx.DB
	{
		sslMode := "require"
		if cfg.DB.DisableTLS {
			sslMode = "disable"
		}
		u := url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(cfg.DB.User, cfg.DB.Password),
			Host:   cfg.DB.Host,
			Path:   cfg.DB.Name,
			RawQuery: (url.Values{
				"sslmode":  []string{sslMode},
				"timezone": []string{"utc"},
			}).Encode(),
		}

		var err error
		db, err = sqlx.Connect("postgres", u.String())
		if err != nil {
			return errors.Wrap(err, "connecting to db")
		}

		defer db.Close()
	}

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db.DB); err != nil {
			return errors.Wrap(err, "applying migrations")
		}
		log.Println("Migrations complete")

	case "seed":
		if err := schema.Seed(db.DB); err != nil {
			return errors.Wrap(err, "seeding database")
		}
		log.Println("Seed data complete")
	}

	return nil
}
