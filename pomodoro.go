package main

import (
	"time"
)

// Pomodoro represents a single unit of time
// to work on a particular task.
type Pomodoro struct {
	// Duration  time.Duration `json:"duration"`
	// End       time.Time     `json:"end"`
	ID        int64         `json:"id"`
	TaskID    int64         `json:"task_id"`
	Start     time.Time     `json:"start"`
	RunTime   time.Duration `json:"run_time"`
	PauseTime time.Duration `json:"pause_time"`
}

// NewPomodoros creates an initialized array of n pomodoro
func NewPomodoros(n int) []*Pomodoro {
	pomodoros := make([]*Pomodoro, n)
	for i := 0; i < n; i++ {
		pomodoros[i] = new(Pomodoro)
	}
	return pomodoros
}
