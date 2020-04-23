package util

import (
	"math/rand"
	"testing"
)

func TestNewRoundRobinIndex(t *testing.T) {
	rr := NewRoundRobinIndex(10)
	var have [21]int
	for i := 0; i < 21; i++ {
		have[i] = rr.Next()
	}
	want := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	if have != want {
		t.Errorf("have:%v, want:%v", have, want)
		return
	}
}

func BenchmarkRoundRobinIndex_Next(b *testing.B) {
	rr := NewRoundRobinIndex(10)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr.Next()
	}
}

func BenchmarkRandIntn(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rand.Intn(10)
	}
}
