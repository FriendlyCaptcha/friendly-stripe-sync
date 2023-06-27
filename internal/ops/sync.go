package ops

import (
	"container/list"
	"context"
	"database/sql"
	"time"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/db/postgres"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/event"
)

func SyncEvents(db *postgres.PostgresStore) error {
	c := context.Background()

	params := &stripe.EventListParams{
		ListParams: stripe.ListParams{
			Limit: stripe.Int64(100),
		},
		Types: []*string{
			stripe.String("customer.created"),
			stripe.String("customer.updated"),
			stripe.String("customer.deleted"),
			stripe.String("product.created"),
			stripe.String("product.updated"),
			stripe.String("product.deleted"),
			stripe.String("price.created"),
			stripe.String("price.updated"),
			stripe.String("price.deleted"),
			stripe.String("subscription.created"),
			stripe.String("subscription.updated"),
			stripe.String("subscription.deleted"),
			stripe.String("coupon.created"),
			stripe.String("coupon.updated"),
			stripe.String("coupon.deleted"),
			stripe.String("customer.discount.updated"),
			stripe.String("customer.discount.deleted"),
		},
	}

	latestSnyc, err := db.Q.GetCurrentSyncState(c)
	if err == nil {
		timeLimit := time.Now().AddDate(0, 0, -30).Unix()
		if latestSnyc.LastEvent < timeLimit {
			log.Warn().Msg("Last sync was more than 28 days ago, do an initial load to make sure there is no missing data")
		}

		params.CreatedRange = &stripe.RangeQueryParams{
			GreaterThan: latestSnyc.LastEvent,
		}
	} else if err == sql.ErrNoRows {
		log.Warn().Msg("No sync state found, you should do an initial load first")
	} else {
		log.Fatal().Err(err).Msg("Failed to get latest sync state")
	}

	log.Info().Int64("last_sync", latestSnyc.LastEvent).Msgf("Starting to load events from stripe")

	events := list.New()
	i := event.List(params)
	for i.Next() {
		e := i.Event()

		// we reverse the list because stripe gives us the events in reverse chronological order
		events.PushFront(e)
	}

	log.Info().Msgf("Finished loading %d events", events.Len())
	log.Info().Msgf("Starting to apply events to database")

	for event := events.Front(); event != nil; event = event.Next() {
		e := event.Value.(*stripe.Event)

		err := HandleEvent(c, db, e)
		if err != nil {
			// if handling an event fails we abort the whole sync because we don't want to miss any events
			log.Error().Err(err).Msg("Failed to handle event")
			return err
		}
		err = db.Q.SetSyncState(c, e.Created)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update sync state")
		}
	}

	log.Info().Msgf("Finished applying all events to database")

	return nil
}
