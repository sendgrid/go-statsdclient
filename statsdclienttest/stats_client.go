package statsdclienttest

import (
	"sync"
	"time"
)

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
	commands map[StatsCommand]int

	// The accumulated values of each stat
	Values map[string]int

	// Whether or not the logger was closed
	Closed bool

	// mutex for making updates to the underlying map atomic
	mutex sync.RWMutex
}

func (m *StatsClient) SetPrefix(prefix string) {

}

func (m *StatsClient) Increment(stat string, delta int, sampleRate float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.commands[StatsCommand{"Increment", stat, delta, sampleRate}] += 1
	m.Values[stat] += delta
	return nil
}

func (m *StatsClient) Decrement(stat string, delta int, sampleRate float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.commands[StatsCommand{"Decrement", stat, delta, sampleRate}] += 1
	m.Values[stat] -= delta
	return nil
}

func (m *StatsClient) Gauge(stat string, value int, sampleRate float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.commands[StatsCommand{"Gauge", stat, value, sampleRate}] += 1
	m.Values[stat] = value
	return nil
}

func (m *StatsClient) Duration(stat string, duration time.Duration, sampleRate float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.commands[StatsCommand{"Duration", stat, int(duration), sampleRate}] += 1
	m.Values[stat] = int(duration / time.Millisecond)
	return nil
}

func (m *StatsClient) Close() error {
	m.Closed = true
	return nil
}

// AssertStat asserts that a given stat is in the list of logged stats, then removes it
func (m *StatsClient) AssertStat(t Testable, stat StatsCommand) {
	m.AssertLoggedN(t, stat, 1)
	delete(m.commands, stat)
}

// AssertValue asserts that the value of the stat with the given stat matches the given value.
// If the stat has not been logged, the test will fail.
func (m *StatsClient) AssertValue(t Testable, stat string, value int) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if _, ok := m.Values[stat]; !ok {
		t.Errorf("expected stat %q to be logged", stat)
	}
}

func (m *StatsClient) AssertLoggedN(t Testable, stat StatsCommand, n int) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.commands[stat] != n {
		t.Errorf("stat %q logged %d times, expected %d times", stat, m.commands[stat], n)
	}
}

// NewStatsClient retur
func NewStatsClient() *StatsClient {
	return &StatsClient{
		commands: make(map[StatsCommand]int),
		Values:   make(map[string]int),
	}
}
