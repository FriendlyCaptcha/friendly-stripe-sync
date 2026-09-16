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

// stripeDocumentedTypesLimit is the cap the Stripe API documents on the event list `types`
// filter: "An array of up to 20 strings containing specific event names." The tests assert
// against this rather than against maxEventTypesPerRequest so the two cannot drift apart.
const stripeDocumentedTypesLimit = 20

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

// TestListEventsSinceChunksRequests drives listEventsSince against a fake Stripe and asserts
// what the 400 from v0.4.0 was really about: no single request may carry more than
// maxEventTypesPerRequest types, and between them the requests must still ask for all of them.
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

		// Two events per request, newest first, so the merged result is only in order if
		// listEventsSince sorts it.
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
		// Deliberately the documented API limit rather than maxEventTypesPerRequest, so that
		// raising our own constant past what Stripe accepts fails here instead of in production.
		assert.LessOrEqual(t, len(types), stripeDocumentedTypesLimit,
			"Stripe rejects a request carrying more than %d types with a 400", stripeDocumentedTypesLimit)
		assert.NotEmpty(t, types)
		all = append(all, types...)
	}

	want := slices.Clone(syncedEventTypes)
	sort.Strings(want)
	sort.Strings(all)
	assert.Equal(t, want, all, "the chunked requests must still ask for every synced type")

	// Four events over two requests, interleaved in time, returned oldest first.
	require.Len(t, events, 4)
	var last int64
	for _, e := range events {
		assert.GreaterOrEqual(t, e.Created, last, "merged events must be applied oldest first")
		last = e.Created
	}
}
