package postgres

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const defaultTargetSessionAttrs = "read-write"

func BuildConnectionDSN(cfg Config) string {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s sslmode=%s connect_timeout=4",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.SSLMode)

	if cfg.Password != "" {
		dsn += " password=" + cfg.Password
	}

	dsn += " target_session_attrs=" + defaultTargetSessionAttrs
	return dsn
}

// validateHostList accepts a single host or a comma-separated list with no whitespace (eg, "10.0.0.1,10.0.0.2").
func validateHostList(host string) error {
	if strings.ContainsFunc(host, unicode.IsSpace) {
		return errors.New("must not contain whitespace; use a host or comma-separated host list")
	}
	for _, h := range strings.Split(host, ",") {
		if h == "" {
			return errors.New("must be a host or comma-separated host list")
		}
	}
	return nil
}

type Store struct {
	db *sqlx.DB
	Q  *Queries
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresStore(cfg Config) (*Store, error) {
	if err := validateHostList(cfg.Host); err != nil {
		return nil, fmt.Errorf("invalid postgres host %q: %w", cfg.Host, err)
	}
	db, err := sqlx.Open("postgres", BuildConnectionDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to postgres store: %w", err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(80)
	db.SetConnMaxLifetime(time.Minute * 15)

	return &Store{
		db: db,
		Q:  New(db),
	}, nil
}

func NewPostgresTestStore(cfg Config) (*Store, func()) {
	// We first need to authenticate against the ordinary dbname
	dbBootstrap, err := sqlx.Open("postgres", BuildConnectionDSN(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to postgres bootstrap store: %v", err)
	}

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	cfgTest := cfg // Copy the config
	cfgTest.DBName = "friendly_stripe_sync_test_" + strings.ToLower(fmt.Sprintf("%X", b))

	_, err = dbBootstrap.Exec(`CREATE DATABASE ` + cfg.DBName + `;`)
	if err != nil {
		log.Fatalf("Couldn't create Postgres test DB: %v", err)
	}

	db, err := sqlx.Open("postgres", BuildConnectionDSN(cfgTest))
	if err != nil {
		log.Fatalf("Failed to connect to postgres test store: %v", err)
	}

	cleanup := func() {
		if err = db.Close(); err != nil {
			fmt.Println("Unable to close Test DB connection: ", err)
			os.Exit(-5)
		}

		_, err = dbBootstrap.Exec(fmt.Sprintf("DROP DATABASE %s", cfgTest.DBName))
		if err != nil {
			log.Fatalf("Couldn't DROP db: %v", err)
		}

		if err = dbBootstrap.Close(); err != nil {
			fmt.Println("Unable to close Bootstrap Test DB connection: ", err)
			os.Exit(-5)
		}

	}

	store := &Store{
		db: db,
		Q:  New(db),
	}
	mig, err := store.GetMigrater(cfgTest)
	if err != nil {
		log.Printf("Failed to get postgres migrater: %v\n", err)
		cleanup()
		os.Exit(1)
	}
	defer mig.Close()

	err = mig.Up()
	if err != nil {
		log.Printf("Failed to migrate Postgres test store: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	return store, cleanup
}
