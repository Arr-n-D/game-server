package server

import (
	"bytes"
	"fmt"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
	"github.com/vmihailenco/msgpack/v5"
)

func (s *Server) networkThread() {
	defer s.threadWaitGroup.Done()

	gns.RunCallbacks()
	for !s.Quit {
		s.pollForIncomingMessages()
	}

}

type Item struct {
	Foo string
}

func (s *Server) pollForIncomingMessages() {
	for !s.Quit {

		messages := make([]*gns.Message, 1)

		mSuccess := s.PollGroup.ReceiveMessages(messages)

		if mSuccess == 0 {
			break
		}

		if mSuccess < 0 {
			sentry.CaptureMessage("Failed to receive messages")
		}
		var item Item
		mPayloadData := messages[0].Payload()

		bytes.NewBuffer(mPayloadData)

		msgpack.Unmarshal(mPayloadData, item)
		fmt.Println(item)
		messages[0].Release()

		// s.ReceiveMessagesChannel <- mPayload
	}
}
