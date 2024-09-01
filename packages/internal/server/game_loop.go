package server

import (
	"log/slog"
	"time"
)

type Item struct {
	Foo string
}

func (s *Server) gameLoopThread() {
	defer s.ThreadWaitGroup.Done()

	lastTickTime := time.Now()
	var tickStartTime time.Time
	var processingTime time.Duration

	for ok := true; ok; ok = !s.Quit {
		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastTickTime)

		if deltaTime >= tickDuration {
			tickStartTime = time.Now()
			// Process stuff here
			s.readIncomingMessages()
			processingTime = time.Since(tickStartTime)
			lastTickTime = currentTime
			slog.With("tick", processingTime).Debug("time to process tick")
		}

		// Yield to other goroutines
		time.Sleep(time.Millisecond)
	}
}

func (s *Server) readIncomingMessages() {
	for {
		select {
		case msg := <-s.ReceiveMessagesChannel:
			// fmt.Println("Case message")
			s.MessagesToProcess = append(s.MessagesToProcess, msg)
			// fmt.Println(len(s.ReceiveMessagesChannel))
		default:
			// fmt.Println("Case default")
			// No more messages available without blocking
			return
		}
	}
}
