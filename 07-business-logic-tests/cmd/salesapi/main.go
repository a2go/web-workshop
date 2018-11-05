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

	"github.com/ardanlabs/service-training/07-business-logic-tests/cmd/salesapi/internal/handlers"
	"github.com/ardanlabs/service-training/07-business-logic-tests/internal/products"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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
		log.Fatal(errors.Wrap(err, "parsing config"))
	}

	if flags.configOnly {
		if err := json.NewEncoder(os.Stdout).Encode(cfg); err != nil {
			log.Fatal(errors.Wrap(err, "encoding config as json"))
		}
		os.Exit(2)
	}

	// Initialize dependencies.
	db, err := sqlx.Connect("postgres", products.DBConn(cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Name, cfg.DB.DisableTLS))
	if err != nil {
		log.Fatal(errors.Wrap(err, "connecting to db"))
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
