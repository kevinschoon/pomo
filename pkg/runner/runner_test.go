package runner_test

import (
	"testing"
	"time"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/runner"
)

func checkState(t *testing.T, state, expected runner.State) {
	t.Helper()
	if state != expected {
		t.Fatalf("expected state %s, got %s", expected, state)
	}
}

func TestTaskRunner(t *testing.T) {
	task := &pomo.Task{
		Duration:  time.Second,
		Message:   "TestTask",
		Pomodoros: pomo.NewPomodoros(4),
	}
	r := runner.New(task)
	r.Start()
	checkState(t, r.State(), runner.INITIALIZED)
	r.Toggle() // start first timer
	checkState(t, r.State(), runner.RUNNING)
	r.Toggle() // noop
	checkState(t, r.State(), runner.RUNNING)
	time.Sleep(time.Second + config.TickTime) // finish first timer
	checkState(t, r.State(), runner.BREAKING)
	r.Toggle() // start second timer
	checkState(t, r.State(), runner.RUNNING)
	time.Sleep(500*time.Millisecond + config.TickTime) // finish half
	r.Suspend()                                        // suspend
	checkState(t, r.State(), runner.SUSPENDED)
	time.Sleep(time.Second + config.TickTime) // suspend one second
	r.Suspend()                               // unsuspend
	checkState(t, r.State(), runner.RUNNING)
	time.Sleep(500*time.Millisecond + config.TickTime) // second half
	checkState(t, r.State(), runner.BREAKING)
	r.Toggle() // third timer
	time.Sleep(time.Second + config.TickTime)
	checkState(t, r.State(), runner.BREAKING)
	r.Toggle() // fourth timer
	time.Sleep(time.Second + config.TickTime)
	checkState(t, r.State(), runner.COMPLETE) // finished
	r.Toggle()                                // shutdown

	for n, pomodoro := range task.Pomodoros {
		t.Logf(
			"Timer %d: start=?,runTime=%s,suspendTime=%s",
			n,
			pomodoro.RunTime,
			pomodoro.PauseTime,
		)
		if pomodoro.RunTime.Truncate(time.Second) != time.Second {
			t.Fatalf(
				"timer %d should have ran for 1s, got %s",
				n,
				pomodoro.RunTime,
			)
		}
	}

	if task.Pomodoros[1].PauseTime.Truncate(time.Second) != time.Second {
		t.Fatalf(
			"second timer should have been suspended 1s, got %s",
			task.Pomodoros[1].PauseTime,
		)
	}
}
