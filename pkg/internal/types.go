package pomo

import (
	"time"

	"github.com/gen2brain/beeep"
)

type State int

func (s State) String() string {
	switch s {
	case RUNNING:
		return "RUNNING"
	case BREAKING:
		return "BREAKING"
	case COMPLETE:
		return "COMPLETE"
	case PAUSED:
		return "PAUSED"
	}
	return ""
}

const (
	RUNNING State = iota + 1
	BREAKING
	COMPLETE
	PAUSED
)

// Wheel keeps track of an ASCII spinner
type Wheel int

func (w *Wheel) String() string {
	switch int(*w) {
	case 0:
		*w++
		return "|"
	case 1:
		*w++
		return "/"
	case 2:
		*w++
		return "-"
	case 3:
		*w = 0
		return "\\"
	}
	return ""
}

// Task describes some activity
type Task struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	// Array of completed pomodoros
	Pomodoros []*Pomodoro `json:"pomodoros"`
	// Free-form tags associated with this task
	Tags []string `json:"tags"`
	// Number of pomodoros for this task
	NPomodoros int `json:"n_pomodoros"`
	// Duration of each pomodoro
	Duration time.Duration `json:"duration"`
}

// ByID is a sortable array of tasks
type ByID []*Task

func (b ByID) Len() int           { return len(b) }
func (b ByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByID) Less(i, j int) bool { return b[i].ID < b[j].ID }

// After returns tasks that were started after the
// provided start time.
func After(start time.Time, tasks []*Task) []*Task {
	filtered := []*Task{}
	for _, task := range tasks {
		if len(task.Pomodoros) > 0 {
			if start.Before(task.Pomodoros[0].Start) {
				filtered = append(filtered, task)
			}
		}
	}
	return filtered
}

// Pomodoro is a unit of time to spend working
// on a single task.
type Pomodoro struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Duration returns the runtime of the pomodoro
func (p Pomodoro) Duration() time.Duration {
	return (p.End.Sub(p.Start))
}

// Status is used to communicate the state
// of a running Pomodoro session
type Status struct {
	State         State         `json:"state"`
	Remaining     time.Duration `json:"remaining"`
	Pauseduration time.Duration `json:"pauseduration"`
	Count         int           `json:"count"`
	NPomodoros    int           `json:"n_pomodoros"`
}

// Returns beeep notifier (no terminal bell) with iconPath set
func NewBeepNotifier(iconPath string) func(string, string) error {
	return func(title string, body string) error {
		return beeep.Notify(title, body, iconPath)
	}
}

// Returns beeep alerter (notifications and terminal bell) with iconPath set
func NewBeepAlerter(iconPath string) func(string, string) error {
	return func(title string, body string) error {
		return beeep.Alert(title, body, iconPath)
	}
}
