package configuration

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func InitSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:                "https://9c07d825e55089ba0cc8a6ad43489f7b@o4507852124323840.ingest.de.sentry.io/4507852128059472",
		ProfilesSampleRate: 1.0,
		TracesSampleRate:   1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)
}
