package server

import (
	"log/slog"
	"time"
)

const (
	tickRate     = 24
	tickDuration = time.Second / tickRate
)

func (s *Server) gameLoopThread() {
	defer s.threadWaitGroup.Done()

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
	select {
	case msg := <-s.ReceiveMessagesChannel:
		s.MessagesToProcess = append(s.MessagesToProcess, msg)
	default:
		return
	}
}
