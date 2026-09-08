package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeConfig writes contents to a temp yaml file and points the config loader at it.
func writeConfig(t *testing.T, contents string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

	oldCfgFile := CfgFile
	CfgFile = path
	t.Cleanup(func() {
		CfgFile = oldCfgFile
		viper.Reset()
	})

	viper.Reset()
	InitConfig()
}

// Keys containing an underscore only bind because we decode using the json tags,
// mapstructure's field-name fallback does not match them.
func TestGetStructUnderscoreKeysRoundTrip(t *testing.T) {
	writeConfig(t, `
stripe:
  api_key: sk_test_fromfile
stripe_sync:
  interval_seconds: 900
  excluded_fields:
    - customer.tax_ids
    - customer.phone
logging:
  filename: /var/log/fss.log
  max_size: 11
  max_age: 22
  max_backups: 33
postgres:
  host: db.example.com
  port: 6543
  dbname: someotherdb
  sslmode: require
`)

	cfg := GetStruct()

	assert.Equal(t, "sk_test_fromfile", cfg.Stripe.APIKey)
	assert.Equal(t, 900, cfg.StripeSync.IntervalSeconds)
	assert.Equal(t, []string{"customer.tax_ids", "customer.phone"}, cfg.StripeSync.ExcludedFields)

	assert.Equal(t, "/var/log/fss.log", cfg.Logging.Filename)
	assert.Equal(t, 11, cfg.Logging.MaxSize)
	assert.Equal(t, 22, cfg.Logging.MaxAge)
	assert.Equal(t, 33, cfg.Logging.MaxBackups)

	// Single-word keys worked before the fix, make sure they still do.
	assert.Equal(t, "db.example.com", cfg.Postgres.Host)
	assert.Equal(t, 6543, cfg.Postgres.Port)
	assert.Equal(t, "someotherdb", cfg.Postgres.DBName)
	assert.Equal(t, "require", cfg.Postgres.SSLMode)
}

func TestGetStructDefaults(t *testing.T) {
	writeConfig(t, "debug: true\n")

	cfg := GetStruct()

	assert.True(t, cfg.Debug)
	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, 5432, cfg.Postgres.Port)
	assert.Equal(t, "friendlystripe", cfg.Postgres.DBName)
	assert.Equal(t, "disable", cfg.Postgres.SSLMode)
}

// Env vars override the config file, including for underscore keys.
func TestGetStructEnvOverride(t *testing.T) {
	t.Setenv("FSS_STRIPE__API_KEY", "sk_test_fromenv")
	t.Setenv("FSS_STRIPE_SYNC__INTERVAL_SECONDS", "42")
	t.Setenv("FSS_POSTGRES__DBNAME", "envdb")

	writeConfig(t, `
stripe:
  api_key: sk_test_fromfile
stripe_sync:
  interval_seconds: 900
`)

	cfg := GetStruct()

	assert.Equal(t, "sk_test_fromenv", cfg.Stripe.APIKey)
	assert.Equal(t, 42, cfg.StripeSync.IntervalSeconds)
	assert.Equal(t, "envdb", cfg.Postgres.DBName)
}

// The underscore-free spellings (stripe.apikey, stripesync.excludedfields, ...) are what
// mapstructure's field-name fallback used to accept, and were the only keys that bound
// before the json tags were honoured. They are deliberately no longer supported.
func TestGetStructLegacyKeysNoLongerBind(t *testing.T) {
	writeConfig(t, `
stripe:
  apikey: sk_test_legacy
stripesync:
  intervalseconds: 77
  excludedfields:
    - customer.address
logging:
  maxsize: 5
`)

	cfg := GetStruct()

	assert.Empty(t, cfg.Stripe.APIKey)
	assert.Zero(t, cfg.StripeSync.IntervalSeconds)
	assert.Empty(t, cfg.StripeSync.ExcludedFields)
	assert.Zero(t, cfg.Logging.MaxSize)
}
