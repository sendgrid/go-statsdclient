package statsdclienttest

import "time"

// Any type used for reporting errors. This will usually be your testing.T variable
// in tests.
type Testable interface {
	Errorf(format string, args ...interface{})
}

// StatsCommand represents a stat that was logged
type StatsCommand struct {
	Operation  string
	Stat       string
	Value      int
	SampleRate float64
}

// A stat logger used for tests
type StatsClient struct {
	// The list of stat commands that have been issued to the stat logger
	Commands []StatsCommand

	// The accumulated values of each stat
	Values map[string]int

	// Whether or not the logger was closed
	Closed bool
}

func (m *StatsClient) Increment(stat string, delta int, sampleRate float64) error {
	m.Commands = append(m.Commands, StatsCommand{"Increment", stat, delta, sampleRate})
	m.Values[stat] += delta
	return nil
}

func (m *StatsClient) Decrement(stat string, delta int, sampleRate float64) error {
	m.Commands = append(m.Commands, StatsCommand{"Decrement", stat, delta, sampleRate})
	m.Values[stat] -= delta
	return nil
}

func (m *StatsClient) Gauge(stat string, value int, sampleRate float64) error {
	m.Commands = append(m.Commands, StatsCommand{"Gauge", stat, value, sampleRate})
	m.Values[stat] = value
	return nil
}

func (m *StatsClient) Timing(stat string, duration time.Duration, sampleRate float64) error {
	m.Commands = append(m.Commands, StatsCommand{"Timing", stat, int(duration), sampleRate})
	m.Values[stat] = int(duration / time.Millisecond)
	return nil
}

func (m *StatsClient) Close() error {
	m.Closed = true
	return nil
}

// AssertStat asserts that a given stat is the first item in the list of commands, then pops that command off
// the list.
func (m *StatsClient) AssertStat(t Testable, stat StatsCommand) {
	actualStat := m.Commands[0]
	if actualStat != stat {
		t.Errorf("got %v stat, expected %v", actualStat, stat)
	}

	// pop it off the top
	m.Commands = m.Commands[1:]
}

// AssertValue asserts that the value of the stat with the given stat matches the given value.
// If the stat has not been logged, the test will fail.
func (m *StatsClient) AssertValue(t Testable, stat string, value int) {
	actualValue, ok := m.Values[stat]
	if !ok {
		t.Errorf("expected stat %q to be logged", stat)
		return
	}
	if actualValue != value {
		t.Errorf("got %d for stat %q, expected %d", actualValue, stat, value)
	}
}

// AssertLogged asserts that any stat with the given stat was logged
func (m *StatsClient) AssertLogged(t Testable, stat string) {
	if _, ok := m.Values[stat]; !ok {
		t.Errorf("expected stat %q to be logged", stat)
	}
}

// NewStatsClient retur
func NewStatsClient() *StatsClient {
	return &StatsClient{
		Commands: make([]StatsCommand, 0),
		Values:   make(map[string]int),
	}
}
