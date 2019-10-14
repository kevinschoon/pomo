package pomo

import (
	"bytes"
	"fmt"
	"time"

	"github.com/bcicen/color"

	"github.com/kevinschoon/pomo/pkg/internal/format"
)

// Pomodoro represents a single unit of time
// to work on a particular task.
type Pomodoro struct {
	// Duration  time.Duration `json:"duration"`
	// End       time.Time     `json:"end"`
	ID        int64         `json:"id,omitempty"`
	TaskID    int64         `json:"task_id,omitempty"`
	Start     time.Time     `json:"start,omitempty"`
	RunTime   time.Duration `json:"run_time,omitempty"`
	PauseTime time.Duration `json:"pause_time,omitempty"`
}

// NewPomodoros creates an initialized array of n pomodoro
func NewPomodoros(n int) []*Pomodoro {
	pomodoros := make([]*Pomodoro, n)
	for i := 0; i < n; i++ {
		pomodoros[i] = new(Pomodoro)
	}
	return pomodoros
}

// Info returns an info string about this Pomodoro
func (p Pomodoro) Info(duration time.Duration) string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "[")
	if p.RunTime >= duration {
		color.New(color.FgHiGreen).Fprintf(buf, "%s", format.TruncDuration(p.RunTime))
	} else {
		color.New(color.FgHiMagenta).Fprintf(buf, "%s", format.TruncDuration(p.RunTime))
	}
	fmt.Fprintf(buf, "]")
	return buf.String()
}
