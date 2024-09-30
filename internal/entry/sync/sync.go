package sync

import (
	"context"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/cfgmodel"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
	"github.com/friendlycaptcha/friendly-stripe-sync/stripesync"
)

func Start(ctx context.Context, cfg cfgmodel.FriendlyStripeSync) error {
	telemetry.SetupLogger(cfg.Development, cfg.Debug, cfg.Logging)

	stripesync, err := stripesync.New(cfg.LibraryConfig())
	if err != nil {
		return fmt.Errorf("failed to create stripe sync: %w", err)
	}

	err = stripesync.SyncEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync events: %w", err)
	}

	return nil
}
