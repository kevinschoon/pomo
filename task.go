package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Task represents a goal to accomplished with
// the Pomodoro technique.
type Task struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Message   string `json:"message"`
	// Array of completed pomodoros
	Pomodoros []*Pomodoro `json:"pomodoros"`
	// Free-form tags associated with this task
	Tags map[string]string `json:"tags"`
	// Number of pomodoros for this task
	// NPomodoros int `json:"n_pomodoros"`
	// Duration of each pomodoro
	Duration time.Duration `json:"duration"`
}

func (t *Task) GetTag(key string) string {
	if t.Tags == nil {
		t.Tags = map[string]string{}
	}
	if value, ok := t.Tags[key]; ok {
		return value
	}
	return ""
}

func (t *Task) SetTag(key, value string) {
	if t.Tags == nil {
		t.Tags = map[string]string{}
	}
	if value, ok := t.Tags[key]; ok {
		if value == "" {
			delete(t.Tags, key)
			return
		}
		t.Tags[key] = value
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
	fmt.Fprintf(buf, " %s", t.Message)
	return buf.String()
}

func truncDuration(s string) string {
	return strings.Replace(strings.Replace(s, "0s", "", -1), "0m", "", -1)
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
