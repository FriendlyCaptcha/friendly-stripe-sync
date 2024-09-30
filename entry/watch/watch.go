package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/friendlycaptcha/friendly-stripe-sync/cfgmodel"
	"github.com/friendlycaptcha/friendly-stripe-sync/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
	"github.com/friendlycaptcha/friendly-stripe-sync/ops"
	"github.com/rs/zerolog/log"
)

func Start(ctx context.Context, cfg cfgmodel.FriendlyStripeSync) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	telemetry.SetupLogger(cfg.Development, cfg.Debug, cfg.Logging)

	intervalSeconds := cfg.StripeSync.IntervalSeconds

	db, err := postgres.NewPostgresStore(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("failed to create postgres store: %w", err)
	}

	stripesync := ops.New(db, cfg.StripeSync, cfg.Stripe.APIKey)

	_, err = db.Q.GetCurrentSyncState(ctx)
	if err == sql.ErrNoRows {
		log.Info().Msg("No sync state found, doing an initial load")
		err := stripesync.InitialLoad(ctx, cfg.Purge)
		if err != nil {
			return fmt.Errorf("failed to load initial data: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to get latest sync state: %w", err)
	}

	log.Info().Msgf("Starting to sync events every %d seconds", intervalSeconds)

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := stripesync.SyncEvents(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Failed to sync events")
			}
		}
	}
}
