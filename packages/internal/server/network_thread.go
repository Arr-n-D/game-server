package server

import (
	"internal/gamemessages"

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

		var msg gamemessages.GameMessage

		// TODO: #14 Grab the information of the messagesPtr like the connection, user data, etc, and put it in message
		decoder := codec.NewDecoderBytes(messagesPtr[0].Payload(), &s.MsgPackHandler)
		err := decoder.Decode(&msg)

		if err != nil {
			sentry.CaptureMessage("An error occured decoding a packet")
			panic("An error occured decoding a packet") // THIS SHOULD NOT BE IN PRODUCTION. DEV ONLY
		}

		msg.MsgNumber = messagesPtr[0].MessageNumber()
		msg.UserData = messagesPtr[0].UserData()
		msg.Connection = messagesPtr[0].Conn()
		msg.ReceivedAt = messagesPtr[0].Timestamp()
		msg.Flags = messagesPtr[0].Flags()

		messagesPtr[0].Release()

		s.ReceiveMessagesChannel <- msg
	}
}
