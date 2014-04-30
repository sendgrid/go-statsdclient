package statsdclient

import (
	"bufio"
	"bytes"
	"testing"
	"time"
)

var result error
var strResult string

func benchFakeClient(buffer *bytes.Buffer) *Client {
	return &Client{
		buf: bufio.NewWriterSize(buffer, defaultBufSize),
	}
}

func BenchmarkIncrement(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)

	for i := 0; i < b.N; i++ {
		r = c.Increment("incr", 1, 1)
	}
	result = r
}

func BenchmarkDecrement(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)

	for i := 0; i < b.N; i++ {
		r = c.Decrement("decr", 1, 1)
	}
	result = r
}

func BenchmarkDuration(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)
	time := time.Duration(123456789)

	for i := 0; i < b.N; i++ {
		r = c.Duration("timing", time, 1)
	}
	result = r
}

func BenchmarkGauge(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)

	for i := 0; i < b.N; i++ {
		r = c.Gauge("gauge", 300, 1)
	}
	result = r
}

func BenchmarkIncrementGauge(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)

	for i := 0; i < b.N; i++ {
		r = c.IncrementGauge("gauge", 10, 1)
	}
	result = r
}

func BenchmarkDecrementGauge(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)

	for i := 0; i < b.N; i++ {
		r = c.DecrementGauge("gauge", 4, 1)
	}
	result = r
}

func BenchmarkUnique(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)

	for i := 0; i < b.N; i++ {
		r = c.Unique("unique", 765, 1)
	}
	result = r
}

func BenchmarkTiming(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)

	for i := 0; i < b.N; i++ {
		r = c.Timing("timing", 350, 1)
	}
	result = r
}

func BenchmarkTime(b *testing.B) {
	var r error
	buf := new(bytes.Buffer)
	c := benchFakeClient(buf)

	for i := 0; i < b.N; i++ {
		r = c.Time("time", 1, func() {})
	}
	result = r
}

func BenchmarkPrefix(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = MakePrefix("test", "statsdclient", "test.example.com")
	}
	strResult = r
}
