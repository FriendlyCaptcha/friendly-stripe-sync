package stripesync_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/friendlycaptcha/friendly-stripe-sync/stripesync"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	syncer, err := stripesync.New(stripesync.Config{
		StripeAPIKey:   "sk_test_dummy",
		ExcludedFields: []string{},
		Postgres: stripesync.PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "friendlystripe",
			SSLMode:  "disable",
		},
	})

	if strings.Contains(err.Error(), "connect: connection refused") {
		t.Skip("Postgres not running or configured, run `docker compose up`")
	}

	require.NoError(t, err)

	err = syncer.MigrateDB(ctx)
	require.NoError(t, err)

	ss, err := syncer.GetCurrentSyncState(ctx)
	assert.ErrorIs(t, err, stripesync.ErrNoSyncState)
	assert.Zero(t, ss)

	err = syncer.InitialLoad(ctx, false)
	// Error should be Stripe invalid request error because of the dummy API key
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid API Key provided")

	err = syncer.SyncEvents(ctx)
	require.ErrorIs(t, err, stripesync.ErrNoSyncState)
}
