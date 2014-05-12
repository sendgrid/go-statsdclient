package statsdclient

import (
	"bufio"
	"bytes"
)

type MockClient struct {
	client
	buffer *bytes.Buffer
}

func (c *MockClient) Buffer() string {
	return c.buffer.String()
}

func (c *MockClient) ResetBuffer() {
	c.buffer.Reset()
}

func (c *MockClient) Close() error {
	return nil
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
//		c.(*statsdclient.MockClient).GetBuffer()
func NewMockClient() *MockClient {
	buffer := new(bytes.Buffer)
	return &MockClient{
		client: client{buf: bufio.NewWriterSize(buffer, defaultBufSize)},
		buffer: buffer,
	}
}
