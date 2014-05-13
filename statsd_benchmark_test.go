package statsdclient

import (
	"testing"
	"time"
)

var result error
var strResult string

func BenchmarkIncrement(b *testing.B) {
	var r error
	c := NewMockClient()

	for i := 0; i < b.N; i++ {
		r = c.Increment("incr", 1, 1)
	}
	// we are assigning r to results in these tests to prevent the compiler
	// from optimizing out the test alltogether (bottom of article)
	// http://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
	result = r
}

func BenchmarkDecrement(b *testing.B) {
	var r error
	c := NewMockClient()

	for i := 0; i < b.N; i++ {
		r = c.Decrement("decr", 1, 1)
	}
	result = r
}

func BenchmarkDuration(b *testing.B) {
	var r error
	c := NewMockClient()
	time := time.Duration(123456789)

	for i := 0; i < b.N; i++ {
		r = c.Duration("timing", time, 1)
	}
	result = r
}

func BenchmarkGauge(b *testing.B) {
	var r error
	c := NewMockClient()

	for i := 0; i < b.N; i++ {
		r = c.Gauge("gauge", 300, 1)
	}
	result = r
}

func BenchmarkIncrementGauge(b *testing.B) {
	var r error
	c := NewMockClient()

	for i := 0; i < b.N; i++ {
		r = c.IncrementGauge("gauge", 10, 1)
	}
	result = r
}

func BenchmarkDecrementGauge(b *testing.B) {
	var r error
	c := NewMockClient()

	for i := 0; i < b.N; i++ {
		r = c.DecrementGauge("gauge", 4, 1)
	}
	result = r
}

func BenchmarkUnique(b *testing.B) {
	var r error
	c := NewMockClient()

	for i := 0; i < b.N; i++ {
		r = c.Unique("unique", 765, 1)
	}
	result = r
}

func BenchmarkTiming(b *testing.B) {
	var r error
	c := NewMockClient()

	for i := 0; i < b.N; i++ {
		r = c.Timing("timing", 350, 1)
	}
	result = r
}

func BenchmarkTime(b *testing.B) {
	var r error
	c := NewMockClient()

	for i := 0; i < b.N; i++ {
		r = c.Time("time", 1, func() {})
	}
	result = r
}
