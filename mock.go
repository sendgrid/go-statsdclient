package statsdclient

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

type MockClient struct {
	client
	buffer *bytes.Buffer
}

func (c *MockClient) Close() error {
	return nil
}

// NextStat returns a string representation of the stat:
// 		Increment "statname:1|c"
// 		Decrement "statname:-1|c"
// 		Duration  "statname:10000.000000|ms"
// 		Gauge     "statname:1|g"
// No newline delimiter is included in the result.
// If no more stats are available, an empty string is returned accompanied by a non-nil error.
func (c *MockClient) NextStat() (string, error) {
	stat, _ := c.buffer.ReadString(0x0A) // newline character
	stat = strings.TrimSpace(stat)

	var err error
	if stat == "" {
		err = errors.New("End of stats")
	}

	return stat, err
}

// Used for mocking the StatsClient for testing purposes
// Using the mock for testing, first wrap the call to Dial in your code appropriately:
// 		var dialStatsd = func(addr string) (StatsClient, error) {
//			return statsdclient.Dial("127.0.0.1:8125")
// 		}
// Then in your test code you can mock out dialStatsd to return the mock object:
// 		dialStatsd = func(addr string) (StatsClient, error) {
//			return statsdclient.NewMockClient(), nil
//		}
// If you want to access the Mock's stats in your test, you'll have to type assert:
//		c.(*statsdclient.MockClient).NextStat()
func NewMockClient() *MockClient {
	buffer := new(bytes.Buffer)
	return &MockClient{
		client: client{buf: bufio.NewWriterSize(buffer, defaultBufSize)},
		buffer: buffer,
	}
}
