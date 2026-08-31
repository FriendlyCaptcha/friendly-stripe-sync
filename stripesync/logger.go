package stripesync

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
)

// IsMissingResourceError reports whether err is a Stripe error for a resource that does not exist.
func IsMissingResourceError(err error) bool {
	var stripeErr *stripe.Error
	return errors.As(err, &stripeErr) && stripeErr.Code == stripe.ErrorCodeResourceMissing
}

type stripeLogger struct {
	downgrade func(error) bool
}

func (stripeLogger) Debugf(format string, v ...interface{}) {
	stripe.DefaultLeveledLogger.Debugf(format, v...)
}

func (stripeLogger) Infof(format string, v ...interface{}) {
	stripe.DefaultLeveledLogger.Infof(format, v...)
}

func (stripeLogger) Warnf(format string, v ...interface{}) {
	stripe.DefaultLeveledLogger.Warnf(format, v...)
}

func (l stripeLogger) Errorf(format string, v ...interface{}) {
	if l.downgrade != nil {
		for _, a := range v {
			if err, ok := a.(error); ok && l.downgrade(err) {
				log.Warn().Msgf(format, v...)
				return
			}
		}
	}

	log.Error().Msgf(format, v...)
}

func stripeBackends(downgrade func(error) bool) *stripe.Backends {
	logger := stripeLogger{downgrade: downgrade}
	newBackend := func(t stripe.SupportedBackend) stripe.Backend {
		return stripe.GetBackendWithConfig(t, &stripe.BackendConfig{LeveledLogger: logger})
	}

	return &stripe.Backends{
		API:     newBackend(stripe.APIBackend),
		Connect: newBackend(stripe.ConnectBackend),
		Uploads: newBackend(stripe.UploadsBackend),
	}
}
