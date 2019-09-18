package runner

import (
	"time"
)

// Status is used to communicate the state
// of a running Pomodoro session
type Status struct {
	Previous      State         `json:"previous"`
	State         State         `json:"state"`
	Count         int           `json:"count"`
	Message       string        `json:"message"`
	NPomodoros    int           `json:"n_pomodoros"`
	Duration      time.Duration `json:"duration"`
	TimeStarted   time.Time     `json:"time_started"`
	TimeRunning   time.Duration `json:"time_running"`
	TimeSuspended time.Duration `json:"time_suspended"`
}
