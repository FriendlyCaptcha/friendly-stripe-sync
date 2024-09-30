package ops

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/cfgmodel"
	"github.com/friendlycaptcha/friendly-stripe-sync/db/postgres"
	"github.com/stripe/stripe-go/v74/client"
)

// Ops is a struct that holds the common global state for all operations.
type Ops struct {
	db     *postgres.PostgresStore
	cfg    cfgmodel.StripeSync
	stripe *client.API
}

// New creates a new Ops struct.
func New(pg *postgres.PostgresStore, cfg cfgmodel.StripeSync, stripeAPIKey string) *Ops {
	stripeClient := &client.API{}
	stripeClient.Init(stripeAPIKey, nil)

	return &Ops{
		db:     pg,
		cfg:    cfg,
		stripe: stripeClient,
	}
}
