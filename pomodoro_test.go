package main

/*

import (
	"testing"
	"time"
)

func TestPomodoro(t *testing.T) {
	// create a new 1 second pomodoro
	pomo := NewPomodoro(1 * time.Second)
	if pomo.Running() {
		t.Fatalf("pomodoro should not be running")
	}
	pomo.Run()
	time.Sleep(100 * time.Millisecond)
	if !pomo.Running() {
		t.Fatalf("pomodoro should be running")
	}
	// toggle (pause) the pomodoro
	pomo.Toggle()
	if pomo.Running() {
		t.Fatalf("pomodoro was paused by still returns running")
	}
	// wait one second
	time.Sleep(1 * time.Second)
	// toggle (unpause) the pomodoro
	pomo.Toggle()
	// wait another 1.1 seconds
	time.Sleep(1100 * time.Millisecond)
	if pomo.Running() {
		t.Fatalf("pomodoro should not be running")
	}
	if !(pomo.Runtime() >= 2*time.Second && pomo.Runtime() < 3*time.Second) {
		t.Fatalf("total duration should be less than 3 but greater than 2 seconds, got %d", pomo.Runtime())
	}
}

*/
