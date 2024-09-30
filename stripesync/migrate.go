package stripesync

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

// MigrateDB automatically migrates the database schema to the latest version.
// This is useful for setting up the database schema for the first time.
//
// This only happens the happy path, if there are any errors, the user will need to fix them manually.
func (o StripeSync) MigrateDB(ctx context.Context) error {
	mig, err := o.db.GetMigrater(o.cfg.Postgres.StoreConfig())
	if err != nil {
		return fmt.Errorf("failed to get migrater: %w", err)
	}

	err = mig.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("failed to migrate up automatically,"+
			" you will need to fix using the friendly-stripe-sync binary (or manual SQL operations): %w", err)
	}

	return nil
}
