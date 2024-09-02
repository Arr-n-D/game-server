package server

import (
	"fmt"

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

		messages := make([]*gns.Message, 1)

		mSuccess := s.PollGroup.ReceiveMessages(messages)

		if mSuccess == 0 {
			break
		}

		if mSuccess < 0 {
			sentry.CaptureMessage("Failed to receive messages")
		}

		var item Item
		var handler codec.MsgpackHandle
		decoder := codec.NewDecoderBytes(messages[0].Payload(), &handler)
		err := decoder.Decode(&item)
		if err != nil {
			panic("Foobar")
		}

		messages[0].Release()
		fmt.Println(item.Foo)
		// mPayloadData := messages[0].Payload()
		// fmt.Println(string(mPayloadData))
		// bytes.NewBuffer(mPayloadData)

		// msgpack.Unmarshal(mPayloadData, item)
		// fmt.Println(item)

		// s.ReceiveMessagesChannel <- mPayload
	}
}
