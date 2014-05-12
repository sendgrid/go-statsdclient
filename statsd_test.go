package statsdclient

import (
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

func TestIncrement(t *testing.T) {
	c := NewMockClient()
	err := c.Increment("incr", 1, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "incr:1|c")
}

func TestDecrement(t *testing.T) {
	c := NewMockClient()
	err := c.Decrement("decr", 1, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "decr:-1|c")
}

func TestDuration(t *testing.T) {
	c := NewMockClient()
	err := c.Duration("timing", time.Duration(123456789), 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "timing:123.456789|ms")
}

func TestIncrementRate(t *testing.T) {
	c := NewMockClient()
	err := c.Increment("incr", 1, 0.99)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "incr:1|c|@0.99")
}

func TestPreciseRate(t *testing.T) {
	c := NewMockClient()
	// The real use case here is rates like 0.0001.
	err := c.Increment("incr", 1, 0.99901)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "incr:1|c|@0.99901")
}

func TestRate(t *testing.T) {
	c := NewMockClient()
	err := c.Increment("incr", 1, 0)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "")
}

func TestGauge(t *testing.T) {
	c := NewMockClient()
	err := c.Gauge("gauge", 300, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "gauge:300|g")
}

func TestIncrementGauge(t *testing.T) {
	c := NewMockClient()
	err := c.IncrementGauge("gauge", 10, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "gauge:+10|g")
}

func TestDecrementGauge(t *testing.T) {
	c := NewMockClient()
	err := c.DecrementGauge("gauge", 4, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "gauge:-4|g")
}

func TestUnique(t *testing.T) {
	c := NewMockClient()
	err := c.Unique("unique", 765, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "unique:765|s")
}

func TestMilliseconds(t *testing.T) {
	msec, _ := time.ParseDuration("350ms")
	assert.Equal(t, 350, millisecond(msec))
	sec, _ := time.ParseDuration("5s")
	assert.Equal(t, 5000, millisecond(sec))
	nsec, _ := time.ParseDuration("50ns")
	assert.Equal(t, 0, millisecond(nsec))
}

func TestTiming(t *testing.T) {
	c := NewMockClient()
	err := c.Timing("timing", 350, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "timing:350|ms")
}

func TestTime(t *testing.T) {
	c := NewMockClient()
	err := c.Time("time", 1, func() { time.Sleep(50e6) })
	assert.Equal(t, err, nil)
}

func TestMultiPacket(t *testing.T) {
	c := NewMockClient()
	err := c.Unique("unique", 765, 1)
	assert.Equal(t, err, nil)
	err = c.Unique("unique", 765, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "unique:765|s\nunique:765|s")
}

func TestMultiPacketOverflow(t *testing.T) {
	c := NewMockClient()
	for i := 0; i < 40; i++ {
		err := c.Unique("unique", 765, 1)
		assert.Equal(t, err, nil)
	}
	assert.Equal(t, c.GetBuffer(), "unique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s\nunique:765|s")
	c.ResetBuffer()
	err := c.Flush()
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetBuffer(), "unique:765|s")
}

func TestPrefix(t *testing.T) {
	c := NewMockClient()
	c.SetPrefix(MakePrefix("test", "statsdclient", "test.example.com"))
	err := c.Increment("key", 1, 1.0)
	assert.Equal(t, err, nil)

	err = c.Flush()
	assert.Equal(t, err, nil)

	assert.Equal(t, c.GetBuffer(), "test.statsdclient.test_example_com.key:1|c")
}
