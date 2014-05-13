package statsdclient

import (
	"bufio"
	"bytes"
	"strings"
)

type MockClient struct {
	client
	buffer *bytes.Buffer
}

func (c *MockClient) Close() error {
	return nil
}

func (c *MockClient) NextStat() string {
	stat, err := c.buffer.ReadString(0x0A) // newline character
	if err == nil {
		stat = strings.TrimSpace(stat)
	}

	return stat
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
// If you want to access the Mock's buffer in your test, you'll have to type assert:
//		c.(*statsdclient.MockClient).Buffer()
func NewMockClient() *MockClient {
	buffer := new(bytes.Buffer)
	return &MockClient{
		client: client{buf: bufio.NewWriterSize(buffer, defaultBufSize)},
		buffer: buffer,
	}
}
