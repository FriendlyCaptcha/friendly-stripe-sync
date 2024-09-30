package sync

import (
	"context"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/cfgmodel"
	"github.com/friendlycaptcha/friendly-stripe-sync/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
	"github.com/friendlycaptcha/friendly-stripe-sync/stripesync"
)

func Start(ctx context.Context, cfg cfgmodel.FriendlyStripeSync) error {
	telemetry.SetupLogger(cfg.Development, cfg.Debug, cfg.Logging)

	db, err := postgres.NewPostgresStore(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("failed to create postgres store: %w", err)
	}

	stripesync := stripesync.New(db, cfg.StripeSync, cfg.Stripe.APIKey)

	err = stripesync.SyncEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync events: %w", err)
	}

	return nil
}
