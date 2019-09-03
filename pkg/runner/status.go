package runner

import (
	"fmt"
	"io"
	"time"
)

// Status is used to communicate the state
// of a running Pomodoro session
type Status struct {
	State         State         `json:"state"`
	Count         int           `json:"count"`
	Message       string        `json:"message"`
	NPomodoros    int           `json:"n_pomodoros"`
	Duration      time.Duration `json:"duration"`
	TimeStarted   time.Time     `json:"time_started"`
	TimeRunning   time.Duration `json:"time_running"`
	TimeSuspended time.Duration `json:"time_suspended"`
}

func (s Status) Write(w io.Writer) {
	state := "?"
	if s.State >= RUNNING {
		state = string(s.State.String()[0])
	}
	if s.State == RUNNING {
		fmt.Fprintf(w, "%s [%d/%d]", state, s.Count, s.NPomodoros)
	} else {
		fmt.Fprintf(w, "%s [%d/%d] -", state, s.Count, s.NPomodoros)
	}
}
