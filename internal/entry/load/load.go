package load

import (
	"github.com/rs/zerolog/log"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/ops"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
)

func Start() {
	telemetry.SetupLogger()

	db := postgres.NewPostgresStore()

	err := ops.InititalLoad(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load initial data")
	}
}
