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
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "incr:1|c")
}

func TestDecrement(t *testing.T) {
	c := NewMockClient()
	err := c.Decrement("decr", 1, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "decr:-1|c")
}

func TestDuration(t *testing.T) {
	c := NewMockClient()
	err := c.Duration("timing", time.Duration(123456789), 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "timing:123.456789|ms")
}

func TestIncrementRate(t *testing.T) {
	c := NewMockClient()
	err := c.Increment("incr", 1, 0.99)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "incr:1|c|@0.99")
}

func TestPreciseRate(t *testing.T) {
	c := NewMockClient()
	// The real use case here is rates like 0.0001.
	err := c.Increment("incr", 1, 0.99901)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "incr:1|c|@0.99901")
}

func TestRate(t *testing.T) {
	c := NewMockClient()
	err := c.Increment("incr", 1, 0)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, err := c.NextStat()
	assert.Equal(t, stat, "")
	assert.NotEqual(t, err, nil)
}

func TestGauge(t *testing.T) {
	c := NewMockClient()
	err := c.Gauge("gauge", 300, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "gauge:300|g")
}

func TestIncrementGauge(t *testing.T) {
	c := NewMockClient()
	err := c.IncrementGauge("gauge", 10, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "gauge:+10|g")
}

func TestDecrementGauge(t *testing.T) {
	c := NewMockClient()
	err := c.DecrementGauge("gauge", 4, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "gauge:-4|g")
}

func TestUnique(t *testing.T) {
	c := NewMockClient()
	err := c.Unique("unique", 765, 1)
	assert.Equal(t, err, nil)
	err = c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "unique:765|s")
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
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "timing:350|ms")
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
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "unique:765|s")
	stat, _ = c.NextStat()
	assert.Equal(t, stat, "unique:765|s")
}

func TestMultiPacketOverflow(t *testing.T) {
	c := NewMockClient()
	for i := 0; i < 40; i++ {
		err := c.Unique("unique", 765, 1)
		assert.Equal(t, err, nil)
	}
	for i := 0; i < 39; i++ {
		stat, _ := c.NextStat()
		assert.Equal(t, stat, "unique:765|s")
	}

	err := c.Flush()
	assert.Equal(t, err, nil)
	stat, _ := c.NextStat()
	assert.Equal(t, stat, "unique:765|s")
}

var prefixTests = []struct {
	prefix   string
	suffix   string
	expected string
}{
	{"test.statsdclient.test_example_com", ".", ".key:1|c"},
	{"test.statsdclient.test_example_com", "", ".key:1|c"},
	{"test.statsdclient.test_example_com", "...", ".key:1|c"},
}

func TestPrefix(t *testing.T) {
	for _, test := range prefixTests {
		c := NewMockClient()

		c.SetPrefix(test.prefix + test.suffix)
		err := c.Increment("key", 1, 1.0)
		assert.Equal(t, err, nil)

		err = c.Flush()
		assert.Equal(t, err, nil)

		stat, _ := c.NextStat()
		assert.Equal(t, stat, test.prefix+test.expected)
	}
}

func TestMultipleCloses(t *testing.T) {
	c := NewMockClient()
	err := c.Close()
	assert.Equal(t, nil, err)
	err = c.Close()
	assert.NotEqual(t, nil, "This should be an error on subsequent closes")
	err = c.Close()
	assert.NotEqual(t, nil, "This should be an error on subsequent closes")
}
