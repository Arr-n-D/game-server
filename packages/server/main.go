package main

import (
	"log/slog"
	"os"
	"time"

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
	if err := server.InitServer(); err != nil {
		slog.Error("failed to server.InitServer()", slog.Any("error", err))
		sentry.CaptureException(err)
		os.Exit(1)
	}

	server.ServerInstance.ThreadWaitGroup.Wait()
	gns.Kill()
	sentry.Flush(5 * time.Second)
}
