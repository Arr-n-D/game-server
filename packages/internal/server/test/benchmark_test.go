package main

import (
	"math/rand"
	"testing"
)

func BenchmarkProcessBatch(b *testing.B) {
	const size = 200
	slice := make([]TickData, size)
	for i := range slice {
		slice[i] = TickData{ID: i, Data: string(rune('A' + i%26))}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processBatch(slice)
	}
}

func processTickData(data *TickData) bool {
	// Simulate processing with random success rate
	return rand.Float32() < 0.8 // 80% success rate
}

func processBatch(slice []TickData) []TickData {
	n := 0
	for i := 0; i < len(slice); i++ {
		if processTickData(&slice[i]) {
			continue
		}
		if i != n {
			slice[n] = slice[i]
		}
		n++
	}
	return slice[:n]
}
