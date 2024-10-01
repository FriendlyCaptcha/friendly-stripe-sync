package watch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/cfgmodel"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
	"github.com/friendlycaptcha/friendly-stripe-sync/stripesync"
	"github.com/rs/zerolog/log"
)

func Start(ctx context.Context, cfg cfgmodel.FriendlyStripeSync) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	telemetry.SetupLogger(cfg.Development, cfg.Debug, cfg.Logging)

	intervalSeconds := cfg.StripeSync.IntervalSeconds

	sync, err := stripesync.New(cfg.LibraryConfig())
	if err != nil {
		return fmt.Errorf("failed to create stripe sync: %w", err)
	}

	_, err = sync.GetCurrentSyncState(ctx)
	if errors.Is(err, stripesync.ErrNoSyncState) {
		log.Info().Msg("No sync state found, doing an initial load")
		err := sync.InitialLoad(ctx, cfg.Purge)
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
			err := sync.SyncEvents(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Failed to sync events")
			}
		}
	}
}
