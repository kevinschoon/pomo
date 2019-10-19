package pomo

import (
	"bytes"
	"fmt"
	"time"

	"github.com/kevinschoon/pomo/pkg/internal/format"
	"github.com/kevinschoon/pomo/pkg/tags"
)

// Task represents a goal to accomplished with
// the Pomodoro technique.
type Task struct {
	ID       int64      `json:"id,omitempty"`
	ParentID int64      `json:"parent_id,omitempty"`
	Message  string     `json:"message,omitempty"`
	Tags     *tags.Tags `json:"tags,omitempty"`
	Tasks    []*Task    `json:"tasks,omitempty"`
	// Array of completed pomodoros
	Pomodoros []*Pomodoro `json:"pomodoros"`
	// Number of pomodoros for this task
	// NPomodoros int `json:"n_pomodoros"`
	// Duration of each pomodoro
	Duration time.Duration `json:"duration"`
}

// NewTask returns a new Task
func NewTask() *Task {
	return &Task{
		Tags: tags.New(),
	}
}

// Info returns an info string about this task
func (t Task) Info() string {
	buf := bytes.NewBuffer(nil)
	pc := int(PercentComplete(t))
	fmt.Fprintf(buf, "[%d]", t.ID)
	fmt.Fprintf(buf, "[%d]", pc)
	fmt.Fprintf(buf, "[%s]", format.TruncDuration(time.Duration(TotalDuration(t))))
	if t.Tags != nil {
		for _, key := range t.Tags.Keys() {
			if t.Tags.Get(key) == "" {
				fmt.Fprintf(buf, "[%s]", key)
			} else {
				fmt.Fprintf(buf, "[%s=%s]", key, t.Tags.Get(key))
			}
		}
	}
	fmt.Fprintf(buf, " %s", t.Message)
	return buf.String()
}
