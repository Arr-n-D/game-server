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
	"internal/messages"

	"github.com/arr-n-d/gns"
	"github.com/ugorji/go/codec"
	"gorm.io/gorm"
)

// TODO: #11 Refactor Server to split things separately. Server shouldn't know about DebugMode or MsgPackHandler or even threadWaitGroup.
type Server struct {
	PollGroup              gns.PollGroup
	listener               *gns.Listener
	Quit                   bool
	threadWaitGroup        sync.WaitGroup
	ReceiveMessagesChannel chan messages.Message
	SendMessagesChannel    chan messages.Message
	MessagesToProcess      []messages.Message
	MsgPackHandler         codec.MsgpackHandle
	DebugMode              bool
	DB                     *gorm.DB
}

func Start(conf *configuration.Configuration) error {
	if conf.Env == configuration.DevEnv || conf.Env == configuration.LocalEnv {
		setDebugOutputFunction(gns.DebugOutputTypeEverything)
	}

	l, err := gns.Listen(&net.UDPAddr{
		IP:   net.IP{127, 0, 0, 1},
		Port: int(conf.GameServerPort),
	},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d. %w", conf.GameServerPort, err)
	}

	poll := gns.NewPollGroup()
	if poll == gns.InvalidPollGroup {
		return errors.New("failed to create poll group")
	}

	database.InitDatabase()

	serverInstance := &Server{
		PollGroup:              poll,
		listener:               l,
		Quit:                   false,
		ReceiveMessagesChannel: make(chan messages.Message, 200),
		SendMessagesChannel:    make(chan messages.Message, 200),
		DebugMode:              false,
		DB:                     database.DATABASE,
	}

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
		slog.Debug("connecting")
		conn := info.Conn()
		if conn.Accept() != gns.ResultOK {
			slog.Error("failed to accept client")
			break
		}

		if !conn.SetPollGroup(s.PollGroup) {
			slog.Error("failed to set poll group")
			break
		}
	}
}
