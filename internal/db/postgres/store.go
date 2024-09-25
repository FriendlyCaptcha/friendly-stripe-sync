package postgres

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config/cfgmodel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func BuildConnectionDSN(cfg cfgmodel.Postgres) string {
	password := cfg.Password

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.SSLMode)

	if password != "" {
		dsn += " password=" + password
	}
	return dsn

}

type PostgresStore struct {
	db *sqlx.DB
	Q  *Queries
}

func NewPostgresStore(cfg cfgmodel.Postgres) *PostgresStore {
	db, err := sqlx.Open("postgres", BuildConnectionDSN(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to postgres store: %v", err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(80)
	db.SetConnMaxLifetime(time.Hour * 1)

	return &PostgresStore{
		db: db,
		Q:  New(db),
	}
}

func NewPostgresTestStore(cfg cfgmodel.Postgres) (*PostgresStore, func()) {
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

	store := &PostgresStore{
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
