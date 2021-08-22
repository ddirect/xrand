package xrand

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestBool(t *testing.T) {
	rnd := New()
	const total = 1e8
	x := 0
	for i := 0; i < total; i++ {
		if rnd.Bool() {
			x++
		}
	}
	prob := float64(x) / total
	err := math.Abs(prob - 0.5)
	limit := 1e-3
	if err > limit {
		t.Errorf("error too high: %f > %f", err, limit)
	}
}

func BenchmarkBool(b *testing.B) {
	rnd := New()
	for i := 0; i < b.N; i++ {
		rnd.Bool()
	}
}

func BenchmarkUniformFactory(b *testing.B) {
	x := New().UniformFactory(0, 0x1000)
	for i := 0; i < b.N; i++ {
		x()
	}
}

func BenchmarkUint64(b *testing.B) {
	rnd := New()
	b.SetBytes(8)
	for i := 0; i < b.N; i++ {
		rnd.Uint64()
	}
}

func BenchmarkUint64NativeLegacy(b *testing.B) {
	rnd := rand.NewSource(time.Now().UnixNano()).(rand.Source64)
	b.SetBytes(8)
	for i := 0; i < b.N; i++ {
		rnd.Uint64()
	}
}
