package server

import (
	"fmt"
	"log"
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
		fmt.Println("Accepted connection")
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
