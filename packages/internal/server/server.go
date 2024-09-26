package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"

	"internal/configuration"
	"internal/database"
	"internal/gamemessages"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
	"github.com/ugorji/go/codec"
	"gorm.io/gorm"
)

// TODO: #11 Refactor Server to split things separately. Server shouldn't know about DebugMode or MsgPackHandler or even threadWaitGroup.
type Server struct {
	PollGroup              gns.PollGroup
	listener               *gns.Listener
	Quit                   bool
	threadWaitGroup        sync.WaitGroup
	ReceiveMessagesChannel chan gamemessages.GameMessage
	SendMessagesChannel    chan gamemessages.GameMessage
	MessagesToProcess      []gamemessages.GameMessage
	MsgPackHandler         codec.MsgpackHandle
	DebugMode              bool
	DB                     *gorm.DB
}

func Start(conf *configuration.Configuration) error {
	if conf.Env == configuration.DevEnv || conf.Env == configuration.LocalEnv {
		setDebugOutputFunction(gns.DebugOutputTypeEverything)
	}

	cfg := gns.ConfigMap{
		gns.ConfigIPAllowWithoutAuth: 1,
	}

	l, err := gns.Listen(&net.UDPAddr{
		IP:   nil,
		Port: int(conf.GameServerPort),
	},
		cfg,
	)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d. %w", conf.GameServerPort, err)
	}

	poll := gns.NewPollGroup()
	if poll == gns.InvalidPollGroup {
		return errors.New("failed to create poll group")
	}

	database.InitDatabase()
	dbInstance := database.GetDatabaseInstance()

	if dbInstance == nil {
		sentry.CaptureMessage("Database instance is nil and shouldn't be")
		panic("Database instance is nil and shouldn't be.")
	}

	serverInstance := &Server{
		PollGroup:              poll,
		listener:               l,
		Quit:                   false,
		ReceiveMessagesChannel: make(chan gamemessages.GameMessage, 200),
		SendMessagesChannel:    make(chan gamemessages.GameMessage, 200),
		DebugMode:              false,
		DB:                     dbInstance,
	}

	serverInstance.InitializeGameFuncsMap()

	dbgMode := os.Getenv("DEBUG")

	if dbgMode == "true" {
		serverInstance.DebugMode = true
	}

	gns.SetGlobalCallbackStatusChanged(serverInstance.StatusCallBackChanged)

	serverInstance.Start()
	return nil
}

func setDebugOutputFunction(detailLevel gns.DebugOutputType) {
	gns.SetDebugOutputFunction(detailLevel, func(typ gns.DebugOutputType, msg string) {
		slog.With("type", typ, "msg", msg).Debug("[DEBUG]")
	})
}

func (s *Server) Start() {
	s.threadWaitGroup.Add(2)
	go s.networkThread()
	go s.gameLoopThread()
	s.threadWaitGroup.Wait()
}

func (s *Server) StatusCallBackChanged(info *gns.StatusChangedCallbackInfo) {
	switch state := info.Info().State(); state {
	case gns.ConnectionStateConnected:
		slog.Debug("accepted connection")
	case gns.ConnectionStateConnecting:
		// slog.Debug("connecting")
		// conn := info.Conn()
		// res := conn.Accept()
		// if res != gns.ResultOK {
		// 	slog.With(res, "response").Error("Failed to accept client with: ")
		// 	// slog.Error("failed to accept client")
		// 	break
		// }

		// if !conn.SetPollGroup(s.PollGroup) {
		// 	slog.Error("failed to set poll group")
		// 	break
		// }
	}
}
