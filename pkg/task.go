package pomo

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/kevinschoon/pomo/pkg/internal/format"
	"github.com/kevinschoon/pomo/pkg/tags"
)

// Task represents a goal to accomplished with
// the Pomodoro technique.
type Task struct {
	ID       int64      `json:"id"`
	ParentID int64      `json:"project_id"`
	Message  string     `json:"message"`
	Tags     *tags.Tags `json:"tags"`
	Tasks    []*Task    `json:"tasks"`
	// Array of completed pomodoros
	Pomodoros []*Pomodoro `json:"pomodoros"`
	// Number of pomodoros for this task
	// NPomodoros int `json:"n_pomodoros"`
	// Duration of each pomodoro
	Duration time.Duration `json:"duration"`
}

func NewTask() *Task {
	return &Task{
		Tags: tags.New(),
	}
}

func (t Task) Info() string {
	buf := bytes.NewBuffer(nil)
	pc := int(PercentComplete(t))
	fmt.Fprintf(buf, "[%d]", t.ID)
	fmt.Fprintf(buf, "[")
	if pc == 100 {
		color.New(color.FgHiGreen).Fprintf(buf, "%d%%", pc)
	} else if pc > 50 && pc < 100 {
		color.New(color.FgHiYellow).Fprintf(buf, "%d%%", pc)
	} else {
		color.New(color.FgHiMagenta).Fprintf(buf, "%d%%", pc)
	}
	fmt.Fprintf(buf, "]")
	fmt.Fprintf(buf, "[%s]", format.TruncDuration(time.Duration(TotalDuration(t))))
	for _, key := range t.Tags.Keys() {
		if t.Tags.Get(key) == "" {
			fmt.Fprintf(buf, "[%s]", key)
		} else {
			fmt.Fprintf(buf, "[%s=%s]", key, t.Tags.Get(key))
		}
	}
	fmt.Fprintf(buf, " %s", t.Message)
	return buf.String()
}

// Fill allocates the remaining time across all
// pomodoros.
func (t *Task) Fill() {
	for _, pomodoro := range t.Pomodoros {
		if pomodoro.Start.IsZero() {
			pomodoro.Start = time.Now()
		}
		pomodoro.RunTime += (t.Duration - pomodoro.RunTime)
	}
}

// Truncate modifies the task so the currently
// allocated time is equal to the duration. Useful
// when a task is completed sooner than expected.
func (t *Task) Truncate() {
	runtime := time.Duration(TimeRunning(*t)).Round(time.Second)
	t.Duration = time.Duration(int64(runtime) / int64(len(t.Pomodoros)))
	for _, pomodoro := range t.Pomodoros {
		pomodoro.Start = time.Time{}
		pomodoro.RunTime = t.Duration
		pomodoro.PauseTime = time.Duration(0)
	}
}
