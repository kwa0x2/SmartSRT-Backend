package bootstrap

import (
	"context"
	"log"
	"log/slog"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	"github.com/kwa0x2/SmartSRT-Backend/config"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}

func InitSentry(env *config.Env) {
	if env.SentryDSN == "" {
		slog.Error("SENTRY_DSN environment variable is empty")
		return
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              env.SentryDSN,
		Environment:      env.AppEnv,
		EnableTracing:    true,
		TracesSampleRate: 0.01,
		EnableLogs:       true,
		Debug:            env.AppEnv == "development",
		MaxBreadcrumbs:   100,
		Transport:        sentry.NewHTTPSyncTransport(),
	}); err != nil {
		slog.Error("Sentry initialization failed", slog.String("error", err.Error()))
		return
	}

	// Console handler için
	var consoleHandler slog.Handler
	if env.AppEnv == "development" {
		consoleHandler = slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	} else {
		consoleHandler = slog.NewJSONHandler(log.Writer(), &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		})
	}

	ctx := context.Background()
	sentryHandler := sentryslog.Option{
		EventLevel: []slog.Level{slog.LevelError},
	}.NewSentryHandler(ctx)

	multiHandler := &MultiHandler{
		handlers: []slog.Handler{consoleHandler, sentryHandler},
	}

	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	slog.Info("Sentry initialized successfully",
		slog.String("environment", env.AppEnv))
}
