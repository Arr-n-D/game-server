package server

import (
	"fmt"
	"internal/configuration"
	"log"
	"net"
	"time"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
)

var ServerInstance *Server

type Server struct {
	PollGroup gns.PollGroup
	listener  *gns.Listener
	Quit      bool

	// Pointer to DB
}

func (s *Server) StatusCallBackChanged(info *gns.StatusChangedCallbackInfo) {
	switch state := info.Info().State(); state {
	case gns.ConnectionStateConnecting:
		fmt.Println("Connecting")
		conn := info.Conn()
		if conn.Accept() != gns.ResultOK {
			log.Fatalln("Failed to accept client")
		}

		if !conn.SetPollGroup(s.PollGroup) {
			log.Fatalln("Failed to set poll group")
		}

	}

}

func InitServer() {
	configuration.InitSentry()
	initGameNetworkingSockets()
	initServer()

	gns.SetGlobalCallbackStatusChanged(ServerInstance.StatusCallBackChanged)
	ServerInstance.pollForIncomingMessages()
	ServerInstance.runCallbacks()
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
		PollGroup: poll,
		listener:  l,
		Quit:      false,
	}

}

func (s *Server) pollForIncomingMessages() {
	go func() {
		for ok := true; ok; ok = !s.Quit {
			fmt.Println("Incoming messages go routine")
			time.Sleep(time.Second * 15)
		}
	}()
}

func (s *Server) runCallbacks() {
	for ok := true; ok; ok = !s.Quit {
		gns.RunCallbacks()
	}
}
