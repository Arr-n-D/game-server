package configuration

import (
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
)

func InitSentry() error {
	env := GetEnv()
	if env == "" {
		return errors.New("no environment found")
	}

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
		return fmt.Errorf("sentry.Init: %w", err)
	}

	return nil
	// sentry.CaptureMessage("Hello!")
}
