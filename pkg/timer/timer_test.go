package timer_test

import (
	"testing"
	"time"

	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/timer"
)

func TestTimer(t *testing.T) {
	runtime := time.Second
	suspendTime := 250 * time.Millisecond
	timer := timer.New(runtime, 0, 0)
	done := timer.Start()
	timer.Suspend()
	time.Sleep(suspendTime)
	timer.Suspend()
	<-done
	maxRunTime := runtime + config.TickTime
	if timer.TimeRunning() < runtime {
		t.Fatalf("should have ran at least %s, got %s", runtime, timer.TimeRunning())
	}
	if timer.TimeRunning() > maxRunTime {
		t.Fatalf("should have ran at most %s, got %s", time.Second+config.TickTime, timer.TimeRunning())
	}
	maxSuspendTime := suspendTime + config.TickTime
	if timer.TimeSuspended() < suspendTime {
		t.Fatalf("should have been suspended at least %s, got %s", suspendTime, timer.TimeSuspended())
	}
	if timer.TimeSuspended() > maxSuspendTime {
		t.Fatalf("should have been suspended at most %s, got %s", maxSuspendTime, timer.TimeSuspended())
	}
	t.Log(timer.TimeRunning(), timer.TimeSuspended())
}
