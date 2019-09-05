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
	ID        int64      `json:"id"`
	ProjectID int64      `json:"project_id"`
	Message   string     `json:"message"`
	Tags      *tags.Tags `json:"tags"`
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
	return running
}

func (t Task) TimePaused() time.Duration {
	var paused time.Duration
	for _, pomodoro := range t.Pomodoros {
		paused += pomodoro.PauseTime
	}
	return paused
}

func (t Task) PercentComplete() float64 {
	duration := int(t.Duration) * len(t.Pomodoros)
	timeRunning := t.TimeRunning()
	return (float64(timeRunning) / float64(duration)) * 100
}

func (t Task) Info() string {
	buf := bytes.NewBuffer(nil)
	pc := int(t.PercentComplete())
	fmt.Fprintf(buf, "[T%d]", t.ID)
	fmt.Fprintf(buf, "[")
	if pc == 100 {
		color.New(color.FgHiGreen).Fprintf(buf, "%d%%", pc)
	} else if pc > 50 && pc < 100 {
		color.New(color.FgHiYellow).Fprintf(buf, "%d%%", pc)
	} else {
		color.New(color.FgHiMagenta).Fprintf(buf, "%d%%", pc)
	}
	fmt.Fprintf(buf, "]")
	fmt.Fprintf(buf, "[%d*%s]", len(t.Pomodoros), format.TruncDuration(t.Duration))
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

// TasksByID is a sortable array of tasks
type TasksByID []*Task

func (b TasksByID) Len() int           { return len(b) }
func (b TasksByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b TasksByID) Less(i, j int) bool { return b[i].ID < b[j].ID }

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
