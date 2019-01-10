package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ardanlabs/service-training/07-logging/cmd/salesapi/internal/handlers"
	"github.com/ardanlabs/service-training/07-logging/internal/platform/log"
	"github.com/ardanlabs/service-training/07-logging/internal/schema"
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

	HTTP struct {
		Address string `default:":8000"`
	}
}

func main() {
	if err := run(); err != nil {
		log.Log("shutting down", "error", err)
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
		log.Log("Migrations complete")
		return nil

	case "seed":
		if err := schema.Seed(db.DB); err != nil {
			return errors.Wrap(err, "seeding database")
		}
		log.Log("Seed data complete")
		return nil
	}

	productsHandler := handlers.Products{DB: db}

	server := http.Server{
		Addr:    cfg.HTTP.Address,
		Handler: http.HandlerFunc(productsHandler.List),
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	log.Log("startup complete")

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "listening and serving")

	case <-osSignals:
		log.Log("caught signal, shutting down")

		// Give outstanding requests 15 seconds to complete.
		const timeout = 15 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Log("gracefully shutting down server", "error", err)
			if err := server.Close(); err != nil {
				log.Log("closing server", "error", err)
			}
		}
	}

	log.Log("done")

	return nil
}
