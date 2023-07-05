package config

import "github.com/spf13/viper"

func setupDefaults() {
	v := viper.GetViper()

	// Postgres defaults
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.dbname", "friendlystripe")
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "")
	v.SetDefault("postgres.sslmode", "disable")
}
