package telemetry

import (
	"fmt"
	"io"
	"os"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/cfgmodel"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(development bool, debug bool, cfg cfgmodel.Logging) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	hostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get hostname")
		hostname = ""
	}

	logContext := log.With().
		Str("host", hostname)

	if debug {
		logContext = logContext.Caller()
	}

	logWriters := make([]io.Writer, 0)
	if development {
		logWriters = append(logWriters, zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		logWriters = append(logWriters, syncWriter())
	}

	if cfg.Filename != "" {
		lj := lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
		}
		logWriters = append(logWriters, &lj)
	}
	writer := io.MultiWriter(logWriters...)
	log.Logger = logContext.Logger().Output(writer)

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// syncWriter is concurrent safe writer.
func syncWriter() io.Writer {
	return diode.NewWriter(os.Stderr, 1000, 0, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	})
}
