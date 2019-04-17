package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/garagesale/internal/platform/auth"
	"github.com/ardanlabs/garagesale/internal/platform/database"
	"github.com/ardanlabs/garagesale/internal/platform/database/schema"
	"github.com/ardanlabs/garagesale/internal/users"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// This is for parsing the environment.
const envKey = "sales"

type config struct {
	DB database.Config
}

func main() {
	if err := run(); err != nil {
		log.Printf("error: %s", err)
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
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "applying migrations")
		}
		log.Println("Migrations complete")

	case "seed":
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "seeding database")
		}
		log.Println("Seed data complete")

	case "useradd":
		email, password := flag.Arg(1), flag.Arg(2)
		if email == "" || password == "" {
			return errors.New("useradd command must be called with two more arguments for email and password")
		}

		fmt.Printf("Admin user will be created with email %q and password %q\n", email, password)
		fmt.Print("Continue? (1/0) ")

		var confirm bool
		if _, err := fmt.Scanf("%t\n", &confirm); err != nil {
			return errors.Wrap(err, "processing response")
		}

		if !confirm {
			fmt.Println("Canceling")
			return nil
		}

		ctx := context.Background()

		// TODO Talk about the engineering decision around validation at web vs service layer.
		nu := users.NewUser{
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Roles:           []string{auth.RoleAdmin, auth.RoleUser},
		}

		u, err := users.Create(ctx, db, nu, time.Now())
		if err != nil {
			return err
		}

		fmt.Printf("User created with id: %v\n", u.ID)
		return nil

	case "keygen":
		path := flag.Arg(1)
		if path == "" {
			return errors.New("keygen missing argument for key path")
		}
		return keygen(path)
	}

	return nil
}

// keygen creates an x509 private key for signing auth tokens.
func keygen(path string) error {

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.Wrap(err, "generating keys")
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "creating private file")
	}
	defer file.Close()

	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if err := pem.Encode(file, &block); err != nil {
		return errors.Wrap(err, "encoding to private file")
	}

	if err := file.Close(); err != nil {
		return errors.Wrap(err, "closing private file")
	}

	return nil
}
