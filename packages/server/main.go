package main

import (
	"log/slog"
	"os"
	"time"

	"internal/configuration"
	"internal/server"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
)

func setupLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	setupLogger()

	conf := configuration.GetConfiguration()
	if conf.Env == "" {
		slog.Error("no configuration found")
		os.Exit(1)
	}

	if err := configuration.InitSentry(); err != nil {
		slog.Error("failed to InitSentry", slog.Any("error", err))
		os.Exit(1)
	}
	defer sentry.Flush(5 * time.Second)

	if err := gns.Init(nil); err != nil {
		slog.Error("failed to initialize gns", slog.Any("error", err))
		os.Exit(1)
	}
	defer gns.Kill()

	if err := server.Start(conf); err != nil {
		slog.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
