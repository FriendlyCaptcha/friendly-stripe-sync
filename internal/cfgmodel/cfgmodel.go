package cfgmodel

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/stripesync"
)

// FriendlyStripeSync is the top-level config for the CLI tool.
type FriendlyStripeSync struct {
	Debug       bool `mapstructure:"debug" json:"debug"`
	Purge       bool `mapstructure:"purge" json:"purge"`
	Development bool `mapstructure:"development" json:"development"`

	Stripe     Stripe     `mapstructure:"stripe" json:"stripe"`
	Postgres   Postgres   `mapstructure:"postgres" json:"postgres"`
	StripeSync StripeSync `mapstructure:"stripe_sync" json:"stripe_sync"`

	Logging Logging `mapstructure:"logging" json:"logging"`
}

type Stripe struct {
	APIKey string `mapstructure:"api_key" json:"api_key"`
}

type Postgres struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	DBName   string `mapstructure:"dbname" json:"dbname"`
	SSLMode  string `mapstructure:"sslmode" json:"sslmode"`
}

type StripeSync struct {
	IntervalSeconds int      `mapstructure:"interval_seconds" json:"interval_seconds"`
	ExcludedFields  []string `mapstructure:"excluded_fields" json:"excluded_fields"`
}

type Logging struct {
	Filename   string `mapstructure:"filename" json:"filename"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
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
