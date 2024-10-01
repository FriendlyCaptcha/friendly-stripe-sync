package stripesync

import (
	"container/list"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
)

// ErrNoSyncState is returned when no sync state is found in the database.
var ErrNoSyncState = fmt.Errorf("no sync state found")

// SyncState contains information about previous sync operations.
type SyncState struct {
	// LastEvent is the timestamp of the last event that was synced in unix time.
	LastEvent int64
	// id is currently always equal to "current_state", and not to be used.
	id string
}

// LastEventTime returns the time of the last event that was synced.
func (s SyncState) LastEventTime() time.Time {
	return time.Unix(s.LastEvent, 0)
}

// MayBeOutdated returns true if the sync state is older than 30 days.
func (s SyncState) MayBeOutdated() bool {
	return time.Since(s.LastEventTime()) > 30*24*time.Hour
}

// SyncEvents syncs all events from stripe (which are up to 30 days old) to the database.
func (o *StripeSync) SyncEvents(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
			stripe.String("customer.subscription.created"),
			stripe.String("customer.subscription.updated"),
			stripe.String("customer.subscription.paused"),
			stripe.String("customer.subscription.deleted"),
			stripe.String("coupon.created"),
			stripe.String("coupon.updated"),
			stripe.String("coupon.deleted"),
			stripe.String("customer.discount.updated"),
			stripe.String("customer.discount.deleted"),
		},
	}

	syncState, err := o.GetCurrentSyncState(ctx)
	if err != nil {
		if !errors.Is(err, ErrNoSyncState) {
			log.Warn().Msg("No sync state found, you should do an initial load first")
		}
		return fmt.Errorf("failed to get latest sync state: %w", err)
	}

	if syncState.MayBeOutdated() {
		log.Warn().Msg("Last sync was more than 30 days ago, do an initial load to make sure there is no missing data")
	}

	log.Info().Int64("last_sync", syncState.LastEvent).Msgf("Starting to load events from stripe")

	events := list.New()
	i := o.stripe.Events.List(params)
	for i.Next() {
		e := i.Event()

		// we reverse the list because stripe gives us the events in reverse chronological order
		events.PushFront(e)
	}
	if err := i.Err(); err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}

	log.Info().Msgf("Finished loading %d events", events.Len())
	log.Info().Msgf("Starting to apply events to database")

	for event := events.Front(); event != nil; event = event.Next() {
		e := event.Value.(*stripe.Event)

		err := o.handleEvent(ctx, e)
		if err != nil {
			// if handling an event fails we abort the whole sync because we don't want to miss any events
			log.Error().Err(err).Msg("Failed to handle event")
			return fmt.Errorf("failed to handle event: %w", err)
		}
		err = o.db.Q.SetSyncState(ctx, e.Created)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update sync state")
		}
	}

	log.Info().Msgf("Finished applying all events to database")

	return nil
}

// GetCurrentSyncState returns the current sync state from the database.
func (o *StripeSync) GetCurrentSyncState(ctx context.Context) (SyncState, error) {
	ss, err := o.db.Q.GetCurrentSyncState(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SyncState{}, ErrNoSyncState
		}
		return SyncState{}, fmt.Errorf("failed to get current sync state: %w", err)
	}

	return SyncState{
		LastEvent: ss.LastEvent,
		id:        ss.ID,
	}, nil
}
