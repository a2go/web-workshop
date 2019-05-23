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
	"github.com/ardanlabs/garagesale/internal/schema"
	"github.com/ardanlabs/garagesale/internal/user"
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

	// This is used for multiple commands below.
	dbConfig := database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	var err error
	switch flag.Arg(0) {
	case "migrate":
		err = migrate(dbConfig)
	case "seed":
		err = seed(dbConfig)
	case "useradd":
		err = useradd(dbConfig, flag.Arg(1), flag.Arg(2))
	case "keygen":
		err = keygen(flag.Arg(1))
	default:
		err = errors.New("Must specify a command")
	}
	if err != nil {
		return err
	}

	return nil
}

func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return err
	}

	fmt.Println("Migrations complete")
	return nil
}

func seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Seed(db); err != nil {
		return err
	}

	fmt.Println("Seed data complete")
	return nil
}

func useradd(cfg database.Config, email, password string) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if email == "" || password == "" {
		return errors.New("useradd command must be called with two additional arguments for email and password")
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

	nu := user.NewUser{
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	u, err := user.Create(ctx, db, nu, time.Now())
	if err != nil {
		return err
	}

	fmt.Println("User created with id:", u.ID)
	return nil
}

// keygen creates an x509 private key for signing auth tokens.
func keygen(path string) error {
	if path == "" {
		return errors.New("keygen missing argument for key path")
	}

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
