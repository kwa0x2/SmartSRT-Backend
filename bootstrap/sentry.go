package bootstrap

import (
	"context"
	"log"
	"log/slog"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	"github.com/kwa0x2/AutoSRT-Backend/config"
)

func InitSentry(env *config.Env) {
	log.Printf("Initializing Sentry with DSN: %s", env.SentryDSN)

	if env.SentryDSN == "" {
		log.Fatalf("SENTRY_DSN environment variable is empty")
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              env.SentryDSN,
		Environment:      env.AppEnv,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		EnableLogs:       true,
		Debug:            env.AppEnv == "development",
		MaxBreadcrumbs:   100,
		Transport:        sentry.NewHTTPSyncTransport(),
	}); err != nil {
		log.Fatalf("Sentry initialization failed: %v", err)
	}

	ctx := context.Background()
	handler := sentryslog.Option{
		EventLevel: []slog.Level{slog.LevelError},
		LogLevel:   []slog.Level{slog.LevelWarn, slog.LevelInfo},
	}.NewSentryHandler(ctx)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	log.Printf("Sentry initialized successfully")
}