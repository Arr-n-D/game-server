package main

import (
	"internal/server"
	"time"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
)

func main() {
	// config.InitSentry()

	server.InitServer()

	server.ServerInstance.ThreadWaitGroup.Wait()
	defer gns.Kill()
	defer sentry.Flush(5 * time.Second)
}
