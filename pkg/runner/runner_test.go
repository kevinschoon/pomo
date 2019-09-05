package runner_test

import (
	"testing"
	"time"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/runner"
)

func TestTaskRunner(t *testing.T) {
	task := &pomo.Task{
		Duration:  time.Second,
		Message:   "TestTask",
		Pomodoros: pomo.NewPomodoros(4),
	}
	var states []runner.State
	r := runner.NewTaskRunner(task, func(status runner.Status) error {
		if len(states) == 0 {
			states = append(states, status.State)
		} else {
			prev := states[len(states)-1]
			if prev != status.State {
				t.Logf("%s -> %s", prev, status.State)
				states = append(states, status.State)
			}
		}
		task.Pomodoros[status.Count].Start = status.TimeStarted
		task.Pomodoros[status.Count].RunTime = status.TimeRunning
		task.Pomodoros[status.Count].PauseTime = status.TimeSuspended
		return nil
	})
	errCh := make(chan error)
	go func() {
		errCh <- r.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	r.Toggle()                                         // start first timer
	r.Toggle()                                         // noop
	time.Sleep(time.Second + config.TickTime)          // finish first timer
	r.Toggle()                                         // start second timer
	time.Sleep(500*time.Millisecond + config.TickTime) // finish half
	r.Suspend()                                        // suspend
	time.Sleep(time.Second + config.TickTime)          // suspend one second
	r.Suspend()                                        // unsuspend
	time.Sleep(500*time.Millisecond + config.TickTime) // second half
	r.Toggle()                                         // third timer
	time.Sleep(time.Second + config.TickTime)
	r.Toggle() // fourth timer
	time.Sleep(time.Second + config.TickTime)
	r.Toggle() // shutdown

	r.Stop()
	err := <-errCh
	if err != nil {
		t.Fatal(err)
	}
	/*
	   runner_test.go:39: INITIALIZED -> RUNNING
	   runner_test.go:39: RUNNING -> BREAKING
	   runner_test.go:39: BREAKING -> RUNNING
	   runner_test.go:39: RUNNING -> SUSPENDED
	   runner_test.go:39: SUSPENDED -> RUNNING
	   runner_test.go:39: RUNNING -> BREAKING
	   runner_test.go:39: BREAKING -> RUNNING
	   runner_test.go:39: RUNNING -> BREAKING
	   runner_test.go:39: BREAKING -> RUNNING
	   runner_test.go:39: RUNNING -> BREAKING
	   runner_test.go:39: BREAKING -> COMPLETE
	*/

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
