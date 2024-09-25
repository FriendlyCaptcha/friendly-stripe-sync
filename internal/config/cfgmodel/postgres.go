package cfgmodel

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
