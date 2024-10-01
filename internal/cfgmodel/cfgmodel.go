package cfgmodel

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/stripesync"
)

// FriendlyStripeSync is the top-level config for the CLI tool.
type FriendlyStripeSync struct {
	Debug       bool `json:"debug"`
	Purge       bool `json:"purge"`
	Development bool `json:"development"`

	Stripe     Stripe     `json:"stripe"`
	Postgres   Postgres   `json:"postgres"`
	StripeSync StripeSync `json:"stripe_sync"`

	Logging Logging `json:"logging"`
}

type Stripe struct {
	APIKey string `json:"api_key"`
}

type Postgres struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type StripeSync struct {
	IntervalSeconds int      `json:"interval_seconds"`
	ExcludedFields  []string `json:"excluded_fields"`
}

type Logging struct {
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

// LibraryConfig returns the config as the stripesync library wants it.
func (c FriendlyStripeSync) LibraryConfig() stripesync.Config {
	return stripesync.Config{
		StripeAPIKey: c.Stripe.APIKey,
		Postgres: stripesync.PostgresConfig{
			Host:     c.Postgres.Host,
			Port:     c.Postgres.Port,
			User:     c.Postgres.User,
			Password: c.Postgres.Password,
			DBName:   c.Postgres.DBName,
			SSLMode:  c.Postgres.SSLMode,
		},
		ExcludedFields: c.StripeSync.ExcludedFields,
	}
}

// PostgresConfig returns the config as the postgres package wants it.
func (c FriendlyStripeSync) PostgresConfig() postgres.Config {
	return postgres.Config{
		Host:     c.Postgres.Host,
		Port:     c.Postgres.Port,
		User:     c.Postgres.User,
		Password: c.Postgres.Password,
		DBName:   c.Postgres.DBName,
		SSLMode:  c.Postgres.SSLMode,
	}
}
