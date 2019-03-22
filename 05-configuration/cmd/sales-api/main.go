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

	"github.com/kelseyhightower/envconfig"

	"github.com/ardanlabs/garagesale/cmd/sales-api/internal/handlers"
	"github.com/ardanlabs/garagesale/internal/platform/database"
)

// This is for parsing the environment.
const envKey = "sales"

type config struct {
	DB   database.Config
	HTTP struct {
		Address         string        `default:":8000"`
		ReadTimeout     time.Duration `default:"5s"`
		WriteTimeout    time.Duration `default:"5s"`
		ShutdownTimeout time.Duration `default:"5s"`
	}
}

func main() {

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
		log.Fatalf("error: parsing config: %s", err)
	}

	// Print config and exit if requested.
	if flags.configOnly {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "	")
		if err := enc.Encode(cfg); err != nil {
			log.Fatalf("error: encoding config as json: %s", err)
		}
		return
	}

	// Initialize dependencies.
	db, err := database.Open(cfg.DB)
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	productsHandler := handlers.Products{DB: db}

	server := http.Server{
		Addr:         cfg.HTTP.Address,
		Handler:      http.HandlerFunc(productsHandler.List),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	log.Println("startup complete")

	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-osSignals:
		log.Println("caught signal, shutting down")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error: gracefully shutting down server: %s", err)
			if err := server.Close(); err != nil {
				log.Printf("error: closing server: %s", err)
			}
		}
	}

	log.Println("done")
}
