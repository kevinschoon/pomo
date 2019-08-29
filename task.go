package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fatih/color"
)

// Task represents a goal to accomplished with
// the Pomodoro technique.
type Task struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Message   string `json:"message"`
	Tags      *Tags  `json:"tags"`
	// Array of completed pomodoros
	Pomodoros []*Pomodoro `json:"pomodoros"`
	// Number of pomodoros for this task
	// NPomodoros int `json:"n_pomodoros"`
	// Duration of each pomodoro
	Duration time.Duration `json:"duration"`
}

func NewTask() *Task {
	return &Task{
		Tags: NewTags(),
	}
}

func (t Task) TimeRunning() time.Duration {
	var d time.Duration
	for _, pomodoro := range t.Pomodoros {
		d += pomodoro.RunTime
	}
	return d
}

func (t Task) TimePaused() time.Duration {
	var d time.Duration
	for _, pomodoro := range t.Pomodoros {
		d += pomodoro.PauseTime
	}
	return d
}

func (t Task) PercentComplete() int {
	completed := t.NCompleted()
	if len(t.Pomodoros) == 0 || completed == 0 {
		return 0
	}
	return (completed / len(t.Pomodoros)) * 100
}

func (t Task) NCompleted() int {
	var n int
	for _, pomodoro := range t.Pomodoros {
		if pomodoro.RunTime >= t.Duration {
			n++
		}
	}
	return n
}

func (t Task) Info() string {
	buf := bytes.NewBuffer(nil)
	pc := t.PercentComplete()
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
	fmt.Fprintf(buf, "[%d*%s]", len(t.Pomodoros), truncDuration(t.Duration.String()))
	for _, key := range t.Tags.Keys() {
		fmt.Fprintf(buf, "[%s=%s]", key, t.Tags.Get(key))
	}
	fmt.Fprintf(buf, " %s", t.Message)
	return buf.String()
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
