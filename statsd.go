/*
Statsd client

Supports counting, sampling, timing, gauges, sets and multi-metrics packet.

Using the client to increment a counter:

	client, err := statsdclient.Dial("127.0.0.1:8125")
	if err != nil {
		// handle error
	}
	defer client.Close()
	err = client.Increment("buckets", 1, 1)

*/
package statsdclient

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultBufSize = 512
)

type StatsClient interface {
	SetPrefix(prefix string)
	Increment(stat string, count int, rate float64) error
	Decrement(stat string, count int, rate float64) error
	Duration(stat string, duration time.Duration, rate float64) error
	Gauge(stat string, value int, rate float64) error
	Close() error
}

// A statsd client representing a connection to a statsd server.
type client struct {
	conn net.Conn
	buf  *bufio.Writer
	m    sync.Mutex

	// The prefix to be added to every key. Should include the "." at the end if desired
	prefix string
}

func millisecond(d time.Duration) int {
	return int(d.Seconds() * 1000)
}

// Dial connects to the given address on the given network using net.Dial and then returns a new client for the connection.
func Dial(addr string) (StatsClient, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	return newClient(conn, 0), nil
}

// DialTimeout acts like Dial but takes a timeout. The timeout includes name resolution, if required.
func DialTimeout(addr string, timeout time.Duration) (StatsClient, error) {
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return nil, err
	}
	return newClient(conn, 0), nil
}

// DialSize acts like Dial but takes a packet size.
// By default, the packet size is 512, see https://github.com/etsy/statsd/blob/master/docs/metric_types.md#multi-metric-packets for guidelines.
func DialSize(addr string, size int) (StatsClient, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	return newClient(conn, size), nil
}

func newClient(conn net.Conn, size int) *client {
	if size <= 0 {
		size = defaultBufSize
	}
	return &client{
		conn: conn,
		buf:  bufio.NewWriterSize(conn, size),
	}
}

// Set the key prefix for the client. All future stats will be sent with the
// prefix value prepended to the bucket.
// Ensures there is only a single "." delimeter at the end. Will remove extraneous ones if present and add one if not present.
func (c *client) SetPrefix(prefix string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.prefix = strings.TrimRight(prefix, ".") + "."
}

// Increment the counter for the given bucket.
func (c *client) Increment(stat string, count int, rate float64) error {
	return c.send(stat, rate, strconv.Itoa(count)+"|c")
}

// Decrement the counter for the given bucket.
func (c *client) Decrement(stat string, count int, rate float64) error {
	return c.Increment(stat, -count, rate)
}

// Record time spent for the given bucket with time.Duration.
func (c *client) Duration(stat string, duration time.Duration, rate float64) error {
	return c.send(stat, rate, strconv.FormatFloat(duration.Seconds()*1000, 'f', 6, 64)+"|ms")
}

// Record time spent for the given bucket in milliseconds.
func (c *client) Timing(stat string, delta int, rate float64) error {
	return c.send(stat, rate, strconv.Itoa(delta)+"|ms")
}

// Calculate time spent in given function and send it.
func (c *client) Time(stat string, rate float64, f func()) error {
	ts := time.Now()
	f()
	return c.Duration(stat, time.Since(ts), rate)
}

// Record arbitrary values for the given bucket.
func (c *client) Gauge(stat string, value int, rate float64) error {
	return c.send(stat, rate, strconv.Itoa(value)+"|g")
}

// Increment the value of the gauge.
func (c *client) IncrementGauge(stat string, value int, rate float64) error {
	return c.send(stat, rate, "+"+strconv.Itoa(value)+"|g")
	// return c.send(stat, rate, "+%d|g", value)
}

// Decrement the value of the gauge.
func (c *client) DecrementGauge(stat string, value int, rate float64) error {
	return c.send(stat, rate, "-"+strconv.Itoa(value)+"|g")
}

// Record unique occurences of events.
func (c *client) Unique(stat string, value int, rate float64) error {
	return c.send(stat, rate, strconv.Itoa(value)+"|s")
}

// Flush writes any buffered data to the network.
func (c *client) Flush() error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.buf.Flush()
}

// Closes the connection.
func (c *client) Close() error {
	c.m.Lock()
	defer c.m.Unlock()
	if c.buf == nil {
		return errors.New("Already closed")
	}
	if err := c.buf.Flush(); err != nil {
		return err
	}
	c.buf = nil
	return c.conn.Close()
}

func (c *client) send(stat string, rate float64, format string, args ...interface{}) error {
	if rate < 1 {
		if rand.Float64() < rate {
			format = format + "|@" + strconv.FormatFloat(rate, 'f', -1, 64)
		} else {
			return nil
		}
	}

	format = c.prefix + stat + ":" + format

	c.m.Lock()
	defer c.m.Unlock()

	// Flush data if we have reach the buffer limit
	if c.buf.Available() < len(format) {
		if err := c.buf.Flush(); err != nil {
			return nil
		}
	}

	// Buffer is not empty, start filling it
	if c.buf.Buffered() > 0 {
		format = "\n" + format
	}

	_, err := fmt.Fprintf(c.buf, format, args...)
	return err
}
