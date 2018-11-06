package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/service-training/08-fatal/cmd/salesapi/internal/handlers"
	"github.com/ardanlabs/service-training/08-fatal/internal/products"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
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
		log.Printf("error: %s", err)
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
	db, err := sqlx.Connect("postgres", products.DBConn(cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Name, cfg.DB.DisableTLS))
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	server := http.Server{
		Addr:    cfg.HTTP.Address,
		Handler: handlers.NewProducts(db),
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
		return errors.Wrap(err, "listening and serving")

	case <-osSignals:
		log.Print("caught signal, shutting down")

		// Give outstanding requests 15 seconds to complete.
		const timeout = 15 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return errors.Wrap(err, "gracefully shutting down server")
		}
	}

	log.Print("done")

	return nil
}
