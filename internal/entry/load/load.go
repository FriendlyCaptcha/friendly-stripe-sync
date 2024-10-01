package load

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

	err = stripesync.InitialLoad(ctx, cfg.Purge)
	if err != nil {
		return fmt.Errorf("failed to load initial data: %w", err)
	}

	return nil
}
