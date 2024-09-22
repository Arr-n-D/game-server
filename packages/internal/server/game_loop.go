package server

import (
	"internal/gamemessages"
	"log/slog"
	"time"

	"github.com/ugorji/go/codec"
)

const (
	tickRate     = 24
	tickDuration = time.Second / tickRate
	maxThreshold = .8
)

// TODO: Read channel of 5000 size. If we're at 80% of threshold, start several Goroutines with batching on the processing
func (s *Server) gameLoopThread() {
	defer s.threadWaitGroup.Done()

	lastTickTime := time.Now()
	var tickStartTime time.Time
	var processingTime time.Duration

	for !s.Quit {
		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastTickTime)

		if deltaTime >= tickDuration {
			tickStartTime = time.Now()
			// Process stuff here
			s.readIncomingMessages(200)
			// slog.With("Size", len(s.MessagesToProcess)).Info("Size of messages is now:")
			s.processBatch()
			processingTime = time.Since(tickStartTime)
			lastTickTime = currentTime
			if s.DebugMode {
				slog.With("tick", processingTime).Debug("time to process tick")
			}
		}

		// Yield to other goroutines
		// time.Sleep(time.Millisecond)
	}
}

func (s *Server) readIncomingMessages(maxMessages int) {
	for i := 0; i < maxMessages; i++ {
		select {
		case msg := <-s.ReceiveMessagesChannel:
			s.MessagesToProcess = append(s.MessagesToProcess, msg)
		default:
			// No more messages in the channel, exit the loop
			return
		}
	}
}

func (s *Server) processTickData(gameMsg *gamemessages.GameMessage) bool {
	msg := gamemessages.MESSAGE_TYPE_TO_TYPE_STRUCT[gameMsg.MessageType]
	decoder := codec.NewDecoderBytes(gameMsg.MessageContent, &s.MsgPackHandler)
	err := decoder.Decode(&msg)
	if err != nil {
		slog.Error("Error decoding message", "error", err)
		return false // Return false if processing failed
	}

	gamemessages.MESSAGE_TYPE_TO_GAME_FUNC[gameMsg.MessageType](msg)
	return true // Return true if processed successfully
}

func (s *Server) processBatch() {
	n := 0
	for i := 0; i < len(s.MessagesToProcess); i++ {
		if !s.processTickData(&s.MessagesToProcess[i]) {
			// If not processed, keep it
			if i != n {
				s.MessagesToProcess[n] = s.MessagesToProcess[i]
			}
			n++
		}
	}
	// Truncate the slice to remove processed messages
	s.MessagesToProcess = s.MessagesToProcess[:n]
}
