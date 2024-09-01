package server

import (
	"log/slog"
	"sync"

	"github.com/arr-n-d/gns"
)

type Server struct {
	PollGroup              gns.PollGroup
	listener               *gns.Listener
	Quit                   bool
	ThreadWaitGroup        sync.WaitGroup
	ReceiveMessagesChannel chan []byte
	SendMessagesChannel    chan []byte
	MessagesToProcess      [][]byte

	// Pointer to DB
}

func (s *Server) init() {
	s.ThreadWaitGroup.Add(1)
	go s.networkThread()
	go s.gameLoopThread()
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
