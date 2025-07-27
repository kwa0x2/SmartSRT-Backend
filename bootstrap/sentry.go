package bootstrap

import (
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kwa0x2/AutoSRT-Backend/config"
)

func InitSentry(env *config.Env) {
	logger := slog.Default()

	environment := func() string {
		if env.AppEnv == "" {
			return "development"
		}
		return env.AppEnv
	}()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              env.SentryDSN,
		Environment:      environment,
		Debug:            env.AppEnv == "development",
		SampleRate:       1.0,
		TracesSampleRate: 0.0,
		Release:          "autosrt-backend@1.0.0",
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			logger.Error("Sentry error captured",
				slog.String("message", event.Message),
				slog.String("event_id", string(event.EventID)),
			)
			return event
		},
	}); err != nil {
		logger.Error("Sentry initialization failed",
			slog.String("error", err.Error()),
			slog.String("environment", environment),
		)
		return
	}

	logger.Info("Sentry initialized successfully",
		slog.String("environment", environment),
		slog.String("status", "initialized"),
	)
}

func CloseSentry() {
	logger := slog.Default()

	if sentry.Flush(2 * time.Second) {
		logger.Info("Sentry closed successfully",
			slog.String("status", "all_events_sent"),
		)
	} else {
		logger.Warn("Sentry closed with timeout",
			slog.String("status", "some_events_may_not_sent"),
		)
	}
}
