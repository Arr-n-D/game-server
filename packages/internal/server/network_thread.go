package server

import (
	"fmt"
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

type Item struct {
	Foo string
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
		}

		var msg messages.Message
		var msg2 messages.Sequence
		var handler codec.MsgpackHandle
		decoder := codec.NewDecoderBytes(messagesPtr[0].Payload(), &handler)
		err := decoder.Decode(&msg)
		if err != nil {
			panic("Foobar")
		}

		fmt.Println(msg.MessageContent)

		decoder = codec.NewDecoderBytes(msg.MessageContent, &handler)
		err = decoder.Decode(&msg2)
		if err != nil {
			panic("Foobar")
		}

		fmt.Println(msg2.Message)
		// mPayloadData := messages[0].Payload()
		// fmt.Println(string(mPayloadData))
		messagesPtr[0].Release()
		// bytes.NewBuffer(mPayloadData)

		// msgpack.Unmarshal(mPayloadData, item)
		// fmt.Println(item)

		// s.ReceiveMessagesChannel <- mPayload
	}
}
