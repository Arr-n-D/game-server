package configuration

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func InitSentry() {
	env := GetEnv()
	Dsn := "https://9c07d825e55089ba0cc8a6ad43489f7b@o4507852124323840.ingest.de.sentry.io/4507852128059472"
	var clientOptions sentry.ClientOptions

	if env == DevEnv {
		clientOptions = sentry.ClientOptions{
			TracesSampleRate:   1.0,
			ProfilesSampleRate: 1.0,
			Dsn:                Dsn,
		}
	} else {
		clientOptions = sentry.ClientOptions{
			TracesSampleRate:   0.2,
			ProfilesSampleRate: 0.2,
			Dsn:                Dsn,
		}
	}

	err := sentry.Init(clientOptions)
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)
}
