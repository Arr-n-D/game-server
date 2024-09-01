package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
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

func (s *Server) networkThread() {

	defer s.ThreadWaitGroup.Done()

	for ok := true; ok; ok = !s.Quit {
		gns.RunCallbacks()
		s.pollForIncomingMessages()
		// s.sendQueuedMessages()
		time.Sleep(time.Millisecond)

	}

}

func (s *Server) pollForIncomingMessages() {

	for ok := true; ok; ok = !s.Quit {

		messages := make([]*gns.Message, 1)

		mSuccess := s.PollGroup.ReceiveMessages(messages)

		if mSuccess == 0 {
			break
		}

		if mSuccess < 0 {
			sentry.CaptureMessage("Failed to receive messages")
		}

		mPayload := make([]byte, len(messages[0].Payload()))
		copy(mPayload, messages[0].Payload())
		messages[0].Release()
		s.ReceiveMessagesChannel <- mPayload
	}
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
