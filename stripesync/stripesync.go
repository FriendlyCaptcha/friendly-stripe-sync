package stripesync

import (
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/stripe/stripe-go/v74/client"
)

// StripeSync is a struct that holds the common global state for all operations.
type StripeSync struct {
	db     *postgres.Store
	stripe *client.API

	cfg Config
}

// Config is the configuration for the StripeSync handle.
type Config struct {
	StripeAPIKey string

	// ExcludedFields is a list of fields that should be excluded from the sync.
	ExcludedFields []string

	Postgres PostgresConfig
}

// PostgresConfig is the configuration for the connection to the Postgres database.
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// New creates a new StripeSync handle.
func New(cfg Config) (*StripeSync, error) {
	stripeClient := &client.API{}
	stripeClient.Init(cfg.StripeAPIKey, nil)

	db, err := postgres.NewPostgresStore(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres store: %w", err)
	}

	return &StripeSync{
		db:     db,
		cfg:    cfg,
		stripe: stripeClient,
	}, nil
}
