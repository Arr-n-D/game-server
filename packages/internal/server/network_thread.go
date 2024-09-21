package server

import (
	"internal/messages"

	"github.com/arr-n-d/gns"
	"github.com/getsentry/sentry-go"
	"github.com/ugorji/go/codec"
)

func (s *Server) networkThread() {
	defer s.threadWaitGroup.Done()

	gns.RunCallbacks()
	for !s.Quit {
		s.pollForIncomingMessages()
	}

}

func (s *Server) pollForIncomingMessages() {
	for !s.Quit {

		messagesPtr := make([]*gns.Message, 1)

		mSuccess := s.PollGroup.ReceiveMessages(messagesPtr)

		if mSuccess == 0 {
			break
		}

		if mSuccess < 0 {
			sentry.CaptureMessage("Failed to receive messages")
			panic("Failed to receive messages")
		}

		var msg messages.Message

		decoder := codec.NewDecoderBytes(messagesPtr[0].Payload(), &s.MsgPackHandler)
		err := decoder.Decode(&msg)
		messagesPtr[0].Release()
		if err != nil {
			// panic("Foobar")
		}

		// decoder = codec.NewDecoderBytes(msg.MessageContent, &handler)
		// err = decoder.Decode(&msg2)
		// if err != nil {
		// 	panic("Foobar")
		// }

		// fmt.Println(msg2.Message)

		s.ReceiveMessagesChannel <- msg
	}
}
