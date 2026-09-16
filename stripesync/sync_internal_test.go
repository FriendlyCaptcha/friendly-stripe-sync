package stripesync

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

// stripeDocumentedTypesLimit is what the Stripe API documents: "An array of up to 20 strings
// containing specific event names." Asserted instead of maxEventTypesPerRequest so raising our
// own constant past it fails here.
const stripeDocumentedTypesLimit = 20

// TestEventTypeChunksCoverEveryType checks that splitting the filter drops no type and reorders
// none. Chunk size is guaranteed by slices.Chunk, so asserting it here would prove nothing.
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

// TestSyncedEventTypesAreHandled keeps the filter and the dispatch in step: a type we request but
// never act on wastes a slot in a filter already at its limit.
func TestSyncedEventTypesAreHandled(t *testing.T) {
	source, err := os.ReadFile("events.go")
	require.NoError(t, err)

	for _, eventType := range syncedEventTypes {
		assert.True(t, strings.Contains(string(source), `"`+eventType+`"`),
			"event type %q is requested from Stripe but has no arm in handleEvent", eventType)
	}
}

// TestSortEventsChronologically covers the ordering the merge depends on: events arrive newest
// first across requests, but the sync state only moves forward if applied oldest first.
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

	// evt_a before evt_b: same second, ID breaks the tie.
	assert.Equal(t, []string{"evt_a", "evt_b", "evt_c", "evt_d"}, got)
}

// TestListEventsSinceChunksRequests drives listEventsSince against a fake Stripe: no request may
// exceed the type limit, and together they must still ask for every type.
func TestListEventsSinceChunksRequests(t *testing.T) {
	var mu sync.Mutex
	var requested [][]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())

		var types []string
		for key, values := range r.Form {
			if strings.HasPrefix(key, "types") {
				types = append(types, values...)
			}
		}
		sort.Strings(types)

		mu.Lock()
		requested = append(requested, types)
		n := len(requested)
		mu.Unlock()

		// Newest first per request, so the merged result is only ordered if listEventsSince sorts.
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprintf(w, `{"object":"list","has_more":false,"url":"/v1/events","data":[
			{"id":"evt_%d_late","object":"event","created":%d,"type":"customer.created","data":{"object":{}}},
			{"id":"evt_%d_early","object":"event","created":%d,"type":"customer.created","data":{"object":{}}}
		]}`, n, 1000+n*10, n, 100+n*10)
		require.NoError(t, err)
	}))
	defer server.Close()

	api := &client.API{}
	api.Init("sk_test_dummy", &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
			URL:           stripe.String(server.URL),
			LeveledLogger: stripe.DefaultLeveledLogger,
		}),
	})

	events, err := (&StripeSync{stripe: api}).listEventsSince(42)
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()

	wantRequests := (len(syncedEventTypes) + maxEventTypesPerRequest - 1) / maxEventTypesPerRequest
	require.Len(t, requested, wantRequests)

	var all []string
	for _, types := range requested {
		// The documented limit, not our constant, so the two cannot drift apart.
		assert.LessOrEqual(t, len(types), stripeDocumentedTypesLimit,
			"Stripe rejects a request carrying more than %d types with a 400", stripeDocumentedTypesLimit)
		assert.NotEmpty(t, types)
		all = append(all, types...)
	}

	want := slices.Clone(syncedEventTypes)
	sort.Strings(want)
	sort.Strings(all)
	assert.Equal(t, want, all, "the chunked requests must still ask for every synced type")

	// Four events across two requests, interleaved in time, returned oldest first.
	require.Len(t, events, 4)
	var last int64
	for _, e := range events {
		assert.GreaterOrEqual(t, e.Created, last, "merged events must be applied oldest first")
		last = e.Created
	}
}
