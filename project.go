package main

import (
	"fmt"
	"io"
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

const (
	emptySpace   = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

type Tree Project

func (t Tree) Write(w io.Writer, depth int, last bool) {
	if depth == 0 { // root
		fmt.Fprintf(w, ".\n")
	}
	spaces := depth
	for i, task := range t.Tasks {
		for j := 0; j < spaces; j++ {
			if j == 0 && !last {
				fmt.Fprintf(w, continueItem)
				continue
			}
			fmt.Fprintf(w, emptySpace)
		}
		if i+1 == len(t.Tasks) && len(t.Children) == 0 {
			fmt.Fprintf(w, lastItem)
		} else {
			fmt.Fprintf(w, middleItem)
		}
		fmt.Fprintf(w, "[T][%d] %s \n", task.ID, task.Message)
	}

	for n, child := range t.Children {
		for j := 0; j < spaces; j++ {
			if j == 0 && !last {
				fmt.Fprintf(w, continueItem)
				continue
			}
			fmt.Fprintf(w, emptySpace)
		}
		if n+1 == len(t.Children) {
			fmt.Fprintf(w, lastItem)
		} else {
			fmt.Fprintf(w, middleItem)
		}
		fmt.Fprintf(w, "[P][%d] %s\n", child.ID, child.Title)
		Tree(*child).Write(w, depth+1, n+1 == len(t.Children) && depth == 0)
	}
}
