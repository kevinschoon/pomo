package main

import (
	"testing"
	"time"
)

func checkState(t *testing.T, state, expected State) {
	t.Helper()
	if state != expected {
		t.Fatalf("expected state %s, got %s", expected, state)
	}
}

func TestTaskRunner(t *testing.T) {
	runner := NewTaskRunner()
	timers := []*Timer{
		NewTimer(time.Second, 0, 0),
		NewTimer(time.Second, 0, 0),
		NewTimer(time.Second, 0, 0),
		NewTimer(time.Second, 0, 0),
	}
	runner.Start(timers)
	checkState(t, runner.State(), INITIALIZED)
	runner.Toggle() // start first timer
	checkState(t, runner.State(), RUNNING)
	runner.Toggle() // noop
	checkState(t, runner.State(), RUNNING)
	time.Sleep(time.Second + tickTime) // finish first timer
	checkState(t, runner.State(), BREAKING)
	runner.Toggle() // start second timer
	checkState(t, runner.State(), RUNNING)
	time.Sleep(500*time.Millisecond + tickTime) // finish half
	runner.Suspend()                            // suspend
	checkState(t, runner.State(), SUSPENDED)
	time.Sleep(time.Second + tickTime) // suspend one second
	runner.Suspend()                   // unsuspend
	checkState(t, runner.State(), RUNNING)
	time.Sleep(500*time.Millisecond + tickTime) // second half
	checkState(t, runner.State(), BREAKING)
	runner.Toggle() // third timer
	time.Sleep(time.Second + tickTime)
	checkState(t, runner.State(), BREAKING)
	runner.Toggle() // fourth timer
	time.Sleep(time.Second + tickTime)
	checkState(t, runner.State(), COMPLETE) // finished
	runner.Toggle()                         // shutdown

	for n, timer := range timers {
		t.Logf(
			"Timer %d: start=?,runTime=%s,suspendTime=%s",
			n,
			timer.TimeRunning(),
			timer.TimeSuspended(),
		)
		if timer.TimeRunning().Truncate(time.Second) != time.Second {
			t.Fatalf(
				"timer %d should have ran for 1s, got %s",
				n,
				timer.TimeRunning(),
			)
		}
	}

	if timers[1].TimeSuspended().Truncate(time.Second) != time.Second {
		t.Fatalf(
			"second timer should have been suspended 1s, got %s",
			timers[1].TimeSuspended(),
		)
	}
}
