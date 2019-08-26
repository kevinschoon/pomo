package main

import (
	"time"
)

// Project is a logical container for tasks
type Project struct {
	ID       int64      `json:"id"`
	ParentID int64      `json:"parent_id"`
	Title    string     `json:"title"`
	Children []*Project `json:"children"`
	Tasks    []*Task    `json:"tasks"`
}

// Duration returns the total length of the project
// including all sub-projects
func (p Project) Duration() time.Duration {
	var duration time.Duration
	for _, child := range p.Children {
		duration += child.Duration()
	}
	for _, task := range p.Tasks {
		duration += task.Duration * time.Duration(int64(len(task.Pomodoros)))
	}
	return duration
}
