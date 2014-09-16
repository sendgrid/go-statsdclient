package statsdclienttest

import (
	"fmt"
	"testing"
	"time"
)

type fakeTestable struct {
	errors []string
}

func (f *fakeTestable) Errorf(format string, args ...interface{}) {
	f.errors = append(f.errors, fmt.Sprintf(format, args...))
}

func (f *fakeTestable) reset() {
	f.errors = make([]string, 0)
}

type assertTestCase struct {
	stats          []StatsCommand
	expectedValues map[string]int
	expectedLogs   []string
	expectedErrors int
}

func TestStatsClient(t *testing.T) {
	testClient := NewStatsClient()

	key := "dec-key"
	delta := 1
	sampleRate := .1
	testClient.Decrement(key, delta, sampleRate)

	expectedCommand := StatsCommand{"Decrement", key, delta, sampleRate}

	testClient.AssertStat(t, expectedCommand)

	key = "inc-key"
	delta = 2
	sampleRate = .2
	testClient.Increment(key, delta, sampleRate)

	expectedCommand = StatsCommand{"Increment", key, delta, sampleRate}
	testClient.AssertStat(t, expectedCommand)

	key = "gauge-key"
	delta = 3
	sampleRate = .3
	testClient.Gauge(key, delta, sampleRate)

	expectedCommand = StatsCommand{"Gauge", key, delta, sampleRate}
	testClient.AssertStat(t, expectedCommand)

	key = "duration-key"
	duration := time.Duration(4) * time.Minute
	sampleRate = .4
	testClient.Duration(key, duration, sampleRate)

	expectedCommand = StatsCommand{"Duration", key, int(duration), sampleRate}
	testClient.AssertStat(t, expectedCommand)
}

func TestAsserts(t *testing.T) {
	testCases := []*assertTestCase{

		// test that increments and decrements work
		&assertTestCase{
			stats: []StatsCommand{
				StatsCommand{
					Operation:  "Increment",
					Stat:       "some-stat",
					Value:      1,
					SampleRate: 1.0,
				},
				StatsCommand{
					Operation:  "Decrement",
					Stat:       "some-stat",
					Value:      2,
					SampleRate: 1.0,
				},
			},
			expectedValues: map[string]int{"some-stat": -1},
			expectedLogs:   []string{"some-stat"},
			expectedErrors: 0,
		},

		// test that gauge sets a value (not increments/decrements)
		&assertTestCase{
			stats: []StatsCommand{
				StatsCommand{
					Operation:  "Increment",
					Stat:       "some-stat",
					Value:      1,
					SampleRate: 1.0,
				},
				StatsCommand{
					Operation:  "Gauge",
					Stat:       "some-stat",
					Value:      5,
					SampleRate: 1.0,
				},
			},
			expectedValues: map[string]int{"some-stat": 5},
			expectedLogs:   []string{"some-stat"},
			expectedErrors: 0,
		},

		// test that errors are reported when assertions fail
		&assertTestCase{
			stats: []StatsCommand{
				StatsCommand{
					Operation:  "Increment",
					Stat:       "some-stat",
					Value:      1,
					SampleRate: 1.0,
				},
			},
			expectedValues: map[string]int{"some-other-stat": 0},
			expectedLogs:   []string{"some-other-stat"},
			expectedErrors: 1,
		},
	}

	for _, tc := range testCases {
		testClient := NewStatsClient()

		// run each stat
		for _, statCmd := range tc.stats {
			switch statCmd.Operation {
			case "Increment":
				testClient.Increment(statCmd.Stat, statCmd.Value, statCmd.SampleRate)
			case "Decrement":
				testClient.Decrement(statCmd.Stat, statCmd.Value, statCmd.SampleRate)
			case "Gauge":
				testClient.Gauge(statCmd.Stat, statCmd.Value, statCmd.SampleRate)
			case "Duration":
				testClient.Duration(statCmd.Stat,
					time.Duration(statCmd.Value)*time.Millisecond,
					statCmd.SampleRate)
			default:
				t.Fatal("unknown operation", statCmd.Operation)
			}
		}

		tester := &fakeTestable{
			errors: make([]string, 0),
		}

		// assert each stat is asserted in AssertStat
		for _, statCmd := range tc.stats {
			testClient.AssertStat(tester, statCmd)

			// always assert 0 for this (instead of len(tc.expectedErorrs))
			// since we iterate against the same set of StatCommands
			// and there should never be an error. we are just checking that
			// AssertStat does not return an error shouldn't, not that it
			// is returning an error when it should.
			if 0 != len(tester.errors) {
				t.Fatalf("AssertStat got %d errors, expected %d:\n%#v",
					len(tester.errors),
					0,
					tester.errors)
			}

			tester.reset()
		}

		for stat, expectedValue := range tc.expectedValues {

			// test AssertValue
			testClient.AssertValue(tester, stat, expectedValue)

			if tc.expectedErrors != len(tester.errors) {
				t.Fatalf("AssertValue got %d errors, expected %d:\n%#v",
					len(tester.errors),
					tc.expectedErrors,
					tester.errors)
			}

			tester.reset()

			// test AssertLogged
			testClient.AssertLogged(tester, stat)

			if tc.expectedErrors != len(tester.errors) {
				t.Fatalf("AssertLogged got %d errors, expected %d:\n%#v",
					len(tester.errors),
					tc.expectedErrors,
					tester.errors)
			}

			tester.reset()
		}
	}
}
