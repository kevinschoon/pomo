package main

import (
	"time"
)

// Status is used to communicate the state
// of a running Pomodoro session
type Status struct {
	State         State         `json:"state"`
	Count         int           `json:"count"`
	NPomodoros    int           `json:"n_pomodoros"`
	TimeStarted   time.Time     `json:"time_started"`
	TimeRunning   time.Duration `json:"time_running"`
	TimeSuspended time.Duration `json:"time_suspended"`
}
