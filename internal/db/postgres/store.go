package postgres

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func BuildConnectionDSN(dbname string) string {
	password := viper.GetString("postgres.password")

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s sslmode=disable",
		viper.GetString("postgres.host"), viper.GetInt("postgres.port"), dbname, viper.GetString("postgres.user"))

	if password != "" {
		dsn += " password=" + password
	}
	return dsn

}

type PostgresStore struct {
	db     *sqlx.DB
	Q      *Queries
	dbname string
}

func NewPostgresStore() *PostgresStore {
	dbname := viper.GetString("postgres.dbname")
	db, err := sqlx.Open("postgres", BuildConnectionDSN(dbname))
	if err != nil {
		log.Fatalf("Failed to connect to postgres store: %v", err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(80)
	db.SetConnMaxLifetime(time.Hour * 1)

	return &PostgresStore{
		db:     db,
		Q:      New(db),
		dbname: dbname,
	}
}

func NewPostgresTestStore() (*PostgresStore, func()) {
	// We first need to authenticate against the ordinary dbname
	dbname := viper.GetString("postgres.dbname")
	dbBootstrap, err := sqlx.Open("postgres", BuildConnectionDSN(dbname))
	if err != nil {
		log.Fatalf("Failed to connect to postgres bootstrap store: %v", err)
	}

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	dbname = "friendly_stripe_sync_test_" + strings.ToLower(fmt.Sprintf("%X", b))
	_, err = dbBootstrap.Exec(`CREATE DATABASE ` + dbname + `;`)
	if err != nil {
		log.Fatalf("Couldn't create Postgres test DB: %v", err)
	}

	db, err := sqlx.Open("postgres", BuildConnectionDSN(dbname))
	if err != nil {
		log.Fatalf("Failed to connect to postgres test store: %v", err)
	}

	cleanup := func() {
		if err = db.Close(); err != nil {
			fmt.Println("Unable to close Test DB connection: ", err)
			os.Exit(-5)
		}

		_, err = dbBootstrap.Exec(fmt.Sprintf("DROP DATABASE %s", dbname))
		if err != nil {
			log.Fatalf("Couldn't DROP db: %v", err)
		}

		if err = dbBootstrap.Close(); err != nil {
			fmt.Println("Unable to close Bootstrap Test DB connection: ", err)
			os.Exit(-5)
		}

	}

	store := &PostgresStore{
		db:     db,
		Q:      New(db),
		dbname: dbname,
	}
	mig, err := store.GetMigrater()
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
