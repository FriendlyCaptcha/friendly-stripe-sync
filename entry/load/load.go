package load

import (
	"context"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/cfgmodel"
	"github.com/friendlycaptcha/friendly-stripe-sync/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
	"github.com/friendlycaptcha/friendly-stripe-sync/ops"
)

func Start(ctx context.Context, cfg cfgmodel.FriendlyStripeSync) error {
	telemetry.SetupLogger(cfg.Development, cfg.Debug, cfg.Logging)

	db, err := postgres.NewPostgresStore(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("failed to create postgres store: %w", err)
	}

	stripesync := ops.New(db, cfg.StripeSync, cfg.Stripe.APIKey)

	err = stripesync.InitialLoad(ctx, cfg.Purge)
	if err != nil {
		return fmt.Errorf("failed to load initial data: %w", err)
	}

	return nil
}
