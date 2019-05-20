package main

import (
	"context"
	"encoding/json"
	_ "expvar" // Register the expvar handlers
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/garagesale/cmd/sales-api/internal/handlers"
	"github.com/ardanlabs/garagesale/internal/platform/database"
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
	HTTP struct {
		Address         string        `default:"localhost:8000"`
		Debug           string        `default:"localhost:6060"`
		ReadTimeout     time.Duration `default:"5s"`
		WriteTimeout    time.Duration `default:"5s"`
		ShutdownTimeout time.Duration `default:"5s"`
	}
}

func main() {
	if err := run(); err != nil {
		log.Println("shutting down", "error:", err)
		os.Exit(1)
	}
}

func run() error {

	log := log.New(os.Stdout, "SALES : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// Process command line flags.
	var flags struct {
		configOnly bool
	}
	flag.BoolVar(&flags.configOnly, "config-only", false, "only show parsed configuration then exit")
	flag.Usage = func() {
		fmt.Print("This program is a service for managing inventory and sales at a Garage Sale.\n\nUsage of sales-api:\n\nsales-api [flags]\n\n")
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

	// =========================================================================
	// Start Debug Service

	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.

	// Not concerned with shutting this down when the application is shutdown.
	go func() {
		log.Println("debug service listening on", cfg.HTTP.Debug)
		err := http.ListenAndServe(cfg.HTTP.Debug, http.DefaultServeMux)
		log.Println("debug service closed", err)
	}()

	// =========================================================================
	// Start API Service

	server := http.Server{
		Addr:         cfg.HTTP.Address,
		Handler:      handlers.API(db, log),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("server listening on", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "listening and serving")

	case <-osSignals:
		log.Println("caught signal, shutting down")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Println("gracefully shutting down server", "error", err)
			if err := server.Close(); err != nil {
				log.Println("closing server", "error", err)
			}
		}
	}

	log.Println("done")

	return nil
}
