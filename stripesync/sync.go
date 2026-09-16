package stripesync

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
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

// maxEventTypesPerRequest is Stripe's cap on the `types` filter. Exceeding it fails the request
// with a 400, so syncedEventTypes is requested in chunks.
const maxEventTypesPerRequest = 20

// syncedEventTypes are the events that change something we mirror. Every entry needs an arm in
// handleEvent.
var syncedEventTypes = []string{
	"customer.created",
	"customer.updated",
	"customer.deleted",
	"customer.tax_id.created",
	"customer.tax_id.updated",
	"customer.tax_id.deleted",
	"product.created",
	"product.updated",
	"product.deleted",
	"price.created",
	"price.updated",
	"price.deleted",
	"customer.subscription.created",
	"customer.subscription.updated",
	"customer.subscription.paused",
	"customer.subscription.deleted",
	"coupon.created",
	"coupon.updated",
	"coupon.deleted",
	"customer.discount.updated",
	"customer.discount.deleted",
}

// SyncEvents syncs all events from stripe (which are up to 30 days old) to the database.
func (o *StripeSync) SyncEvents(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	events, err := o.listEventsSince(syncState.LastEvent)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		log.Info().Msg("Finished loading events, no new events found")
		return nil
	}

	log.Info().Msgf("Finished loading %d events, starting to apply events to database", len(events))

	// Skip `coupon.created` and `coupon.updated` events that have a `coupon.deleted` event later on for the same ID.
	// Sometimes we create/update and quickly delete a coupon. The problem is that Stripe will return a 404
	// if we try to retrieve information about a deleted coupon.
	newOrUpdatedCoupons := make(map[string]bool)
	skipCoupons := make(map[string]bool)
	for _, e := range events {
		if e.Type == "coupon.created" || e.Type == "coupon.updated" {
			id, ok := e.Data.Object["id"].(string)
			if ok {
				newOrUpdatedCoupons[id] = true
			}
		} else if e.Type == "coupon.deleted" {
			id, ok := e.Data.Object["id"].(string)
			if ok && newOrUpdatedCoupons[id] {
				skipCoupons[id] = true
			}
		}
	}

	for _, e := range events {
		if e.Type == "coupon.created" || e.Type == "coupon.updated" {
			id, ok := e.Data.Object["id"].(string)
			if ok && skipCoupons[id] {
				log.Info().Str("coupon_id", id).Msg("Skipping coupon that was deleted")
				continue
			}
		}

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

// listEventsSince returns every event of a synced type created after since, oldest first.
//
// Chunk boundaries are harmless: the merged events are sorted before any is applied, and a
// handler reaching an entity another chunk has not loaded fetches it from Stripe itself.
func (o *StripeSync) listEventsSince(since int64) ([]*stripe.Event, error) {
	var events []*stripe.Event

	for chunk := range slices.Chunk(syncedEventTypes, maxEventTypesPerRequest) {
		i := o.stripe.Events.List(&stripe.EventListParams{
			ListParams:   stripe.ListParams{Limit: stripe.Int64(100)},
			CreatedRange: &stripe.RangeQueryParams{GreaterThan: since},
			Types:        stripe.StringSlice(chunk),
		})
		for i.Next() {
			events = append(events, i.Event())
		}
		if err := i.Err(); err != nil {
			return nil, fmt.Errorf("failed to list events: %w", err)
		}
	}

	sortEventsChronologically(events)

	return events, nil
}

// sortEventsChronologically orders events oldest first. The sync state advances to each event's
// created time as it is applied, so it only moves forward in this order. The ID breaks ties
// within a second so a batch always applies the same way.
func sortEventsChronologically(events []*stripe.Event) {
	slices.SortStableFunc(events, func(a, b *stripe.Event) int {
		if c := cmp.Compare(a.Created, b.Created); c != 0 {
			return c
		}
		return cmp.Compare(a.ID, b.ID)
	})
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
