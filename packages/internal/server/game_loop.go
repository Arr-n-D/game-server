package server

import (
	"fmt"
	"time"
)

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
			s.processIncomingMessages()
			processingTime = time.Since(tickStartTime)
			lastTickTime = currentTime
			fmt.Printf("Time to process tick: %v\n", processingTime)
		}

		// Yield to other goroutines
		time.Sleep(time.Millisecond)
	}

}

func (s *Server) processIncomingMessages() {
    
}
