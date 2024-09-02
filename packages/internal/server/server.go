package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"internal/configuration"
	"internal/messages"

	"github.com/arr-n-d/gns"
)

type Server struct {
	PollGroup              gns.PollGroup
	listener               *gns.Listener
	Quit                   bool
	threadWaitGroup        sync.WaitGroup
	ReceiveMessagesChannel chan messages.Message
	SendMessagesChannel    chan messages.Message
	MessagesToProcess      [][]byte
	// MsgPackHandler         codec.MsgpackHandle
	// Pointer to DB
}

func Start(conf *configuration.Configuration) error {
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
		return fmt.Errorf("failed to listen on port %d. %w", conf.GameServerPort, err)
	}

	poll := gns.NewPollGroup()
	if poll == gns.InvalidPollGroup {
		return errors.New("failed to create poll group")
	}

	serverInstance := &Server{
		PollGroup:              poll,
		listener:               l,
		Quit:                   false,
		ReceiveMessagesChannel: make(chan messages.Message, 200),
		SendMessagesChannel:    make(chan messages.Message, 200),
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
		}

		if !conn.SetPollGroup(s.PollGroup) {
			slog.Error("failed to set poll group")
		}
	}
}
