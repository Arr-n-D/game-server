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

		// TODO: #14 Grab the information of the messagesPtr like the connection, user data, etc, and put it in message
		decoder := codec.NewDecoderBytes(messagesPtr[0].Payload(), &s.MsgPackHandler)
		err := decoder.Decode(&msg)

		if err != nil {
			sentry.CaptureMessage("An error occured decoding a packet")
			panic("An error occured decoding a packet") // THIS SHOULD NOT BE IN PRODUCTION. DEV ONLY
		}

		messagesPtr[0].Release()

		// decoder = codec.NewDecoderBytes(msg.MessageContent, &handler)
		// err = decoder.Decode(&msg2)
		// if err != nil {
		// 	panic("Foobar")
		// }

		// fmt.Println(msg2.Message)

		s.ReceiveMessagesChannel <- msg
	}
}
