package main

import (
	"context"
	"fmt"

	"github.com/friendlycaptcha/friendly-stripe-sync/stripesync"
)

// This is an example of how to use the stripesync library.

func main() {
	cfg := stripesync.Config{
		StripeAPIKey: "<stripe-api-key>",
		ExcludedFields: []string{
			"customer.address",
		},
		Postgres: stripesync.PostgresConfig{
			// ...
		},
	}

	stripesync, err := stripesync.New(cfg)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	err = stripesync.MigrateDB(ctx)
	if err != nil {
		panic(err)
	}

	ss, err := stripesync.GetCurrentSyncState(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The timestamp of the last synced Stripe event is %s.\n", ss.LastEventTime())
	fmt.Println("Syncing events...")

	err = stripesync.SyncEvents(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Events synced successfully.")
}
