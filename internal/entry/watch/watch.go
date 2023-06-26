package watch

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/ops"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Start() {
	telemetry.SetupLogger()

	intervalSeconds := viper.GetInt("stripe_sync.interval_seconds")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	db := postgres.NewPostgresStore()

	_, err := db.Q.GetCurrentSyncState(context.Background())
	if err == sql.ErrNoRows {
		log.Info().Msg("No sync state found, doing an initial load")
		err := ops.InititalLoad(db)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to load initial data")
		}
	} else if err != nil {
		log.Fatal().Err(err).Msg("Failed to get latest sync state")
	}

	log.Info().Msgf("Starting to sync events every %d seconds", intervalSeconds)

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	for {
		select {
		case <-sc:
			ticker.Stop()
			return
		case <-ticker.C:
			err := ops.SyncEvents(db)
			if err != nil {
				log.Error().Err(err).Msg("Failed to sync events")
			}
		}
	}
}
