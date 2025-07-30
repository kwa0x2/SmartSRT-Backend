package bootstrap

import (
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/kwa0x2/AutoSRT-Backend/config"
)

func InitSentry(env *config.Env) {

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              env.SentryDSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Fatalf("Sentry initialization failed: %v", err)
	}
}

