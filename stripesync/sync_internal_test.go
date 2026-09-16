package stripesync

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v74"
)

// TestEventTypeChunksCoverEveryType is what keeps chunking honest. The per-chunk size is
// guaranteed by slices.Chunk, so the thing worth asserting is that splitting the filter neither
// drops a type nor reorders one.
func TestEventTypeChunksCoverEveryType(t *testing.T) {
	var covered []string
	chunks := 0
	for chunk := range slices.Chunk(syncedEventTypes, maxEventTypesPerRequest) {
		require.NotEmpty(t, chunk)
		chunks++
		covered = append(covered, chunk...)
	}

	assert.Equal(t, syncedEventTypes, covered, "chunking must not drop or reorder a type")
	assert.Equal(t, (len(syncedEventTypes)+maxEventTypesPerRequest-1)/maxEventTypesPerRequest, chunks)
}

// TestSyncedEventTypesAreHandled keeps the filter and the dispatch in step. A type we ask Stripe
// for but never act on costs a request slot and misleads the next person to count them.
func TestSyncedEventTypesAreHandled(t *testing.T) {
	source, err := os.ReadFile("events.go")
	require.NoError(t, err)

	for _, eventType := range syncedEventTypes {
		assert.True(t, strings.Contains(string(source), `"`+eventType+`"`),
			"event type %q is requested from Stripe but has no arm in handleEvent", eventType)
	}
}

// TestSortEventsChronologically covers the ordering the merged chunks depend on: events arrive
// newest first and interleaved across requests, but the sync state only moves forward if they
// are applied oldest first.
func TestSortEventsChronologically(t *testing.T) {
	events := []*stripe.Event{
		{ID: "evt_d", Created: 300},
		{ID: "evt_b", Created: 100},
		{ID: "evt_c", Created: 200},
		{ID: "evt_a", Created: 100},
	}

	sortEventsChronologically(events)

	got := make([]string, 0, len(events))
	var last int64
	for _, e := range events {
		assert.GreaterOrEqual(t, e.Created, last, "events must never step backwards in time")
		last = e.Created
		got = append(got, e.ID)
	}

	// evt_a before evt_b: same second, so the ID breaks the tie deterministically.
	assert.Equal(t, []string{"evt_a", "evt_b", "evt_c", "evt_d"}, got)
}
