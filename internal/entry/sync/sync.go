package sync

import (
	"context"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config/cfgmodel"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/ops"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
)

func Start(ctx context.Context, cfg cfgmodel.FriendlyStripeSync) error {
	telemetry.SetupLogger(cfg.Development, cfg.Debug, cfg.Logging)

	db := postgres.NewPostgresStore(cfg.Postgres)

	stripesync := ops.New(db, cfg.StripeSync, cfg.Stripe.APIKey)

	err := stripesync.SyncEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync events: %w", err)
	}

	return nil
}
