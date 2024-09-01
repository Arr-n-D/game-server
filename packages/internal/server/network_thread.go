package server

import (
	"time"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
)

func (s *Server) networkThread() {
	defer s.threadWaitGroup.Done()

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
