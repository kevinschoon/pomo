package runner

import (
	"fmt"
	"time"

	"github.com/kevinschoon/pomo/pkg/internal/format"
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

func (s Status) String() string {
	state := "?"
	if s.State >= RUNNING {
		state = string(s.State.String()[0])
	}
	remaining := format.TruncDuration(s.Duration - s.TimeRunning)
	return fmt.Sprintf("%s [%d/%d] - %s", state, s.Count, s.NPomodoros, remaining)
}
