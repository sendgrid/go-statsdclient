package statsdclient

import "time"

type nullStatsClient struct{}

func (n *nullStatsClient) SetPrefix(prefix string) {

}

func (n *nullStatsClient) Increment(stat string, count int, rate float64) error {
	return nil
}

func (n *nullStatsClient) Decrement(stat string, count int, rate float64) error {
	return nil
}

func (n *nullStatsClient) Duration(stat string, duration time.Duration, rate float64) error {
	return nil
}

func (n *nullStatsClient) Gauge(stat string, value int, rate float64) error {
	return nil
}

func (n *nullStatsClient) Close() error {
	return nil
}

// NullStatsClient is a statsdclient that does nothing. This can
// be used in testing.
var NullStatsClient StatsClient = &nullStatsClient{}
