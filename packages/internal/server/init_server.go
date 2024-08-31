package server

import (
	"internal/configuration"
	"log"
	"net"
	"time"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
)

var ServerInstance *Server

const (
	tickRate     = 24
	tickDuration = time.Second / tickRate
)

func InitServer() {
	configuration.InitSentry()
	initGameNetworkingSockets()
	initServer()

	gns.SetGlobalCallbackStatusChanged(ServerInstance.StatusCallBackChanged)

	ServerInstance.init()
}

func initGameNetworkingSockets() {
	err := gns.Init(nil)

	if err != nil {
		log.Fatal(err)
		sentry.CaptureException(err)
	}

}

func setDebugOutputFunction(detailLevel gns.DebugOutputType) {
	gns.SetDebugOutputFunction(detailLevel, func(typ gns.DebugOutputType, msg string) {
		log.Print("[DEBUG] ", typ, msg)
	})
}

func initServer() {
	conf := configuration.GetConfiguration()

	if conf.Env == configuration.DevEnv {
		setDebugOutputFunction(gns.DebugOutputTypeEverything)
	}

	l, err := gns.Listen(&net.UDPAddr{
		IP:   net.IP{127, 0, 0, 1},
		Port: int(conf.GameServerPort),
	},
		nil,
	)

	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Failed to listen on port %d", conf.GameServerPort)
	}

	poll := gns.NewPollGroup()
	if poll == gns.InvalidPollGroup {
		sentry.CaptureMessage("Failed to create poll group")
		log.Fatal("Invalid poll group")
	}

	ServerInstance = &Server{
		PollGroup:              poll,
		listener:               l,
		Quit:                   false,
		ReceiveMessagesChannel: make(chan []byte, 200),
		SendMessagesChannel:    make(chan []byte, 200),
	}

}
