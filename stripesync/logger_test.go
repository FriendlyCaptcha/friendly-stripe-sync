package stripesync

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/coupon"
)

func TestStripeLoggerErrorLevel(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		downgrade func(error) bool
		want      string
	}{
		{"missing resource without downgrade", `{"error":{"code":"resource_missing","type":"invalid_request_error"}}`, nil, `"level":"error"`},
		{"missing resource with downgrade", `{"error":{"code":"resource_missing","type":"invalid_request_error"}}`, IsMissingResourceError, `"level":"warn"`},
		{"other error with downgrade", `{"error":{"code":"api_key_expired","type":"invalid_request_error"}}`, IsMissingResourceError, `"level":"error"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			var buf bytes.Buffer
			prev := log.Logger
			log.Logger = zerolog.New(&buf)
			defer func() { log.Logger = prev }()

			b := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
				LeveledLogger: stripeLogger{downgrade: tt.downgrade},
				URL:           stripe.String(srv.URL),
			})
			client := coupon.Client{B: b, Key: "sk_test_dummy"}
			_, _ = client.Get("co_missing", nil)

			assert.Contains(t, buf.String(), tt.want)
		})
	}
}
