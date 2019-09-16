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

func (t Task) TimeRunning() time.Duration {
	var running time.Duration
	for _, pomodoro := range t.Pomodoros {
		running += pomodoro.RunTime
	}
	for _, subTask := range t.Tasks {
		running += subTask.TimeRunning()
	}
	return running
}

func (t Task) TimePaused() time.Duration {
	var paused time.Duration
	for _, pomodoro := range t.Pomodoros {
		paused += pomodoro.PauseTime
	}
	for _, subTask := range t.Tasks {
		paused += subTask.TimePaused()
	}
	return paused
}

func (t Task) TotalDuration() time.Duration {
	duration := int(t.Duration) * len(t.Pomodoros)
	for _, subTask := range t.Tasks {
		duration += int(subTask.TotalDuration())
	}
	return time.Duration(duration)
}

func (t Task) PercentComplete() float64 {
	if t.TotalDuration() == 0 {
		return 100
	}
	duration := t.TotalDuration()
	timeRunning := t.TimeRunning()
	return (float64(timeRunning) / float64(duration)) * 100
}

func (t Task) Info() string {
	buf := bytes.NewBuffer(nil)
	pc := int(t.PercentComplete())
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
	fmt.Fprintf(buf, "[%s]", format.TruncDuration(t.TotalDuration()))
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
	runtime := t.TimeRunning().Round(time.Second)
	t.Duration = time.Duration(int64(runtime) / int64(len(t.Pomodoros)))
	for _, pomodoro := range t.Pomodoros {
		pomodoro.Start = time.Time{}
		pomodoro.RunTime = t.Duration
		pomodoro.PauseTime = time.Duration(0)
	}
}
