package sync

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/ops"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
	"github.com/rs/zerolog/log"
)

func Start() {
	telemetry.SetupLogger()

	db := postgres.NewPostgresStore()

	err := ops.SyncEvents(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to sync events")
	}
}
