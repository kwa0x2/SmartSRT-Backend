package bootstrap

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kwa0x2/AutoSRT-Backend/config"
)

func InitSentry(env *config.Env) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: env.SentryDSN,
		Environment: func() string {
			if env.AppEnv == "" {
				return "development"
			}
			return env.AppEnv
		}(),
		Debug:            env.AppEnv == "development",
		SampleRate:       1.0,
		TracesSampleRate: 1.0,
		Release:          "autosrt-backend@1.0.0",
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			log.Printf("🚨 Sentry Error: %s", event.Message)
			return event
		},
	}); err != nil {
		log.Fatalf("Sentry initialization failed: %v", err)
	}

	log.Printf("Sentry initialized successfully with environment: %s", func() string {
		if env.AppEnv == "" {
			return "development"
		}
		return env.AppEnv
	}())
}

func CloseSentry() {
	if sentry.Flush(2 * time.Second) {
		log.Println("Sentry closed successfully - all events sent")
	} else {
		log.Println("Sentry closed with timeout - some events may not have been sent")
	}
}
