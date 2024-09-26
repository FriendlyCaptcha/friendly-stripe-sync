package migrate

import (
	"context"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config/cfgmodel"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/store"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/telemetry"
	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type MigrateOpts struct {
	TargetVersion int
}

func Migrate(ctx context.Context, cfg cfgmodel.FriendlyStripeSync, storeName string, operation string, opts MigrateOpts) error {
	// By setting `development` to true we make sure we flush the logs to stdout always.
	cfg.Development = true

	telemetry.SetupLogger(cfg.Development, cfg.Debug, cfg.Logging)

	// Contextual logger
	l := log.With().Str("entry", "migrate").Str("store", storeName).Str("operation", operation).Logger()
	l.Debug().Msg("Starting migration")

	var migrater store.Migrater

	switch storeName {
	case "postgres":
		pg, err := postgres.NewPostgresStore(cfg.Postgres)
		if err != nil {
			return fmt.Errorf("failed to create postgres store: %w", err)
		}
		pgMigrater, err := pg.GetMigrater(cfg.Postgres)
		if err != nil {
			return fmt.Errorf("failed to get migrater: %w", err)
		}
		migrater = pgMigrater
		defer migrater.Close()
	default:
		l.WithLevel(zerolog.FatalLevel).Msg("Unknown store, can't migrate")
		return fmt.Errorf("unknown store, can't migrate: %s", storeName)
	}

	migrater.SetLogger(migrationZeroLogger{
		zerologger: l,
		verbose:    cfg.Debug,
	})

	var err error
	switch operation {
	case "up":
		err = migrater.Up()
	case "down":
		err = migrater.Down()
	case "list":
		var migrations []string
		migrations, err = migrater.List()
		if err != nil {
			break
		}
		l.Info().Strs("migrations", migrations).Msg("")
	case "version":
		var version uint
		var dirty bool
		version, dirty, err = migrater.Version()
		if err != nil {
			break
		}
		l.Info().Uint("version", version).Bool("dirty", dirty).Msg("")

	case "force":
		l = l.With().Int("target_version", opts.TargetVersion).Logger()
		err = migrater.Force(opts.TargetVersion)
		if err != nil {
			break
		}

	case "to":
		l = l.With().Int("target_version", opts.TargetVersion).Logger()
		if opts.TargetVersion < 0 {
			l.WithLevel(zerolog.FatalLevel).Err(err).Msg("Invalid target version for migrate")
		}
		err = migrater.To(uint(opts.TargetVersion))
		if err != nil {
			break
		}
	}

	if err == migrate.ErrNoChange {
		l.Warn().Msg("Already at the correct version, migration was skipped")
	} else if err == migrate.ErrNilVersion {
		l.Warn().Msg("Migration is at nil version (no migrations have been performed)")
	} else if err != nil {
		l.WithLevel(zerolog.FatalLevel).Err(err).Msg("Migration operation failed")
		return fmt.Errorf("migration operation failed: %w", err)
	}

	l.Debug().Msg("Migration end")

	return err
}
