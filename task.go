package main

import (
	"time"
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
	Tags []string `json:"tags"`
	// Number of pomodoros for this task
	// NPomodoros int `json:"n_pomodoros"`
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
