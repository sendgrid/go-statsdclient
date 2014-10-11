package statsdclient

import (
	"io"
	"io/ioutil"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/bmizerany/assert"
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

func TestUDPServerCloses(t *testing.T) {
	addr := net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// create a UDP listener
	listener, err := net.ListenUDP("udp", &addr)
	if err != nil {
		t.Fatal(err)
	}

	waitForCloseChan := make(chan struct{})

	go func() {
		_, err := io.Copy(ioutil.Discard, listener)
		if err != nil {
			t.Log(err)
			close(waitForCloseChan)
		}
	}()

	// Dial to whatever UDP server was created
	conn, err := Dial(listener.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}

	// We need the client for the flush method - the StatsClient interface does not expose it
	client, ok := conn.(*client)
	if !ok {
		t.Fatal("do not have a client - it is needed to manually trigger a flush")
	}

	err = client.Gauge(strings.Repeat("k.", 256), 1, 1)
	if err != nil {
		t.Fatal("could not send Gauge", err)
	}

	err = client.Flush()
	if err != nil {
		t.Fatal("could not flush client", err)
	}

	err = listener.Close()
	if err != nil {
		t.Fatal("could not close listener", err)
	}

	// wait until the listener is fully shut down
	<-waitForCloseChan

	// This loop is required, because for some reason it would not fail consistently on the first attempt
	// before the change to use `exconn` as the conncection object
	// if at first you don't, keep trying and hope you don't
	for i := 0; i < 10; i++ {
		err = client.Gauge(strings.Repeat("k.", 256), 1, 1.0)
		if err != nil {
			t.Fatalf("Attempt #%d: Error should not have occurred even if the UDP server is down: %s", i, err)
		}

		err = client.Flush()
		if err != nil {
			t.Fatalf("Attempt #%d: could not flush client: %s", i, err)
		}
	}
}
