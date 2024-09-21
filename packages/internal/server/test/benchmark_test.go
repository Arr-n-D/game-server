package main

// func BenchmarkProcessBatch(b *testing.B) {
// 	const size = 200
// 	slice := make([]TickData, size)
// 	for i := range slice {
// 		slice[i] = TickData{ID: i, Data: string(rune('A' + i%26))}
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		processBatch(slice)
// 	}
// }

// func (s *Server) readIncomingMessages(maxMessages int) {
// 	for i := 0; i < maxMessages; i++ {
// 		select {
// 		case msg := <-s.ReceiveMessagesChannel:
// 			s.MessagesToProcess = append(s.MessagesToProcess, msg)
// 		default:
// 			// No more messages in the channel, exit the loop
// 			return
// 		}
// 	}
// }

// func (s *Server) processTickData(message *messages.Message) bool {
// 	var msg messages.Sequence
// 	decoder := codec.NewDecoderBytes(message.MessageContent, &s.MsgPackHandler)
// 	err := decoder.Decode(&msg)
// 	if err != nil {
// 		slog.Error("Error decoding message", "error", err)
// 		return false // Return false if processing failed
// 	}
// 	slog.Info("Processing message", "message", msg.Message)
// 	// Add your processing logic here
// 	return true // Return true if processed successfully
// }

// func (s *Server) processBatch() {
// 	n := 0
// 	for i := 0; i < len(s.MessagesToProcess); i++ {
// 		if !s.processTickData(&s.MessagesToProcess[i]) {
// 			// If not processed, keep it
// 			if i != n {
// 				s.MessagesToProcess[n] = s.MessagesToProcess[i]
// 			}
// 			n++
// 		}
// 	}
// 	// Truncate the slice to remove processed messages
// 	s.MessagesToProcess = s.MessagesToProcess[:n]
// }
