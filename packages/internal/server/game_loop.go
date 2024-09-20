package server

import (
	"fmt"
	"internal/messages"
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

	for !s.Quit {
		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastTickTime)

		if deltaTime >= tickDuration {
			tickStartTime = time.Now()
			// Process stuff here
			s.readIncomingMessages()
			processBatch(s.MessagesToProcess)
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

func (s *Server) readIncomingMessages() {
	select {
	case msg := <-s.ReceiveMessagesChannel:
		s.MessagesToProcess = append(s.MessagesToProcess, msg)
	default:
		return
	}
}

func processTickData(message *messages.Message) bool {
	// Simulate processing
	fmt.Printf("Processing: %+v\n", *message)
	return true // Return true if processed successfully
}

// TODO: #13 Two processing lanes, One for reliable and one for unreliable?
func processBatch(slice []messages.Message) []messages.Message {
	n := 0
	for i := 0; i < len(slice); i++ {
		if processTickData(&slice[i]) {
			// Element processed, skip it
			continue
		}
		// Element not processed, keep it
		if i != n {
			slice[n] = slice[i]
		}
		n++
	}
	return slice[:n]
}
