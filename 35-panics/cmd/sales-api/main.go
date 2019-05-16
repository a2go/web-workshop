package main

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	_ "expvar" // Register the expvar handlers
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/garagesale/cmd/sales-api/internal/handlers"
	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/database"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"
)

// This is for parsing the environment.
const envKey = "sales"

type config struct {
	DB struct {
		User     string `default:"postgres"`
		Password string `default:"postgres" json:"-"` // Prevent the marshalling of secrets.
		Host     string `default:"localhost"`
		Name     string `default:"postgres"`
	}
	HTTP struct {
		Address         string        `default:"localhost:8000"`
		Debug           string        `default:"localhost:6060"`
		ReadTimeout     time.Duration `default:"5s"`
		WriteTimeout    time.Duration `default:"5s"`
		ShutdownTimeout time.Duration `default:"5s"`
	}
	Auth struct {
		KeyID          string `default:"1" split_words:"true"`
		PrivateKeyFile string `default:"private.pem" split_words:"true"`
		Algorithm      string `default:"RS256"`
	}
	Trace struct {
		URL         string  `default:"http://localhost:9411/api/v2/spans"`
		Service     string  `default:"sales-api"`
		Probability float64 `default:"1"`
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
	authenticator, err := createAuth(cfg)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}

	db, err := database.Open(database.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Host:     cfg.DB.Host,
		Name:     cfg.DB.Name,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	closer, err := registerTracer(cfg)
	if err != nil {
		return err
	}
	defer closer()

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
		Handler:      handlers.API(db, log, authenticator),
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

func createAuth(cfg config) (*auth.Authenticator, error) {

	keyContents, err := ioutil.ReadFile(cfg.Auth.PrivateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	public := auth.NewSingleKeyFunc(cfg.Auth.KeyID, key.Public().(*rsa.PublicKey))

	return auth.NewAuthenticator(key, cfg.Auth.KeyID, cfg.Auth.Algorithm, public)
}

func registerTracer(cfg config) (func() error, error) {
	localEndpoint, err := openzipkin.NewEndpoint(cfg.Trace.Service, cfg.HTTP.Address)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}
	reporter := zipkinHTTP.NewReporter(cfg.Trace.URL)

	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndpoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(cfg.Trace.Probability),
	})

	return reporter.Close, nil
}
