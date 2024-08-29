package main

import (
	"fmt"
	"internal/server"
	"time"

	"log"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
)

func StatusCallBackChanged(info *gns.StatusChangedCallbackInfo) {
	switch state := info.Info().State(); state {
	case gns.ConnectionStateConnecting:
		fmt.Println("Connecting")
		conn := info.Conn()
		if conn.Accept() != gns.ResultOK {
			log.Fatalln("Failed to accept client")
		}

		if !conn.SetPollGroup(server.ServerInstance.PollGroup) {
			log.Fatalln("Failed to set poll group")
		}

	}

}

func main() {
	// config.InitSentry()

	server.InitServer()
	

	defer gns.Kill()
	defer sentry.Flush(5 * time.Second)
}
