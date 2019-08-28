package main

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
)

// Project is a logical container for tasks
type Project struct {
	ID       int64      `json:"id"`
	ParentID int64      `json:"parent_id"`
	Title    string     `json:"title"`
	Children []*Project `json:"children"`
	Tasks    []*Task    `json:"tasks"`
}

type ProjectFn func(Project)

func ForEach(p Project, fn ProjectFn) {
	fn(p)
	for _, child := range p.Children {
		ForEach(*child, fn)
	}
}

func Duration(p Project) time.Duration {
	var duration time.Duration
	ForEach(p, func(p Project) {
		for _, task := range p.Tasks {
			duration += task.Duration * time.Duration(int64(len(task.Pomodoros)))
		}
	})
	return duration
}

func PercentComplete(p Project) int {
	var i, j int
	ForEach(p, func(p Project) {
		for _, task := range p.Tasks {
			i += len(task.Pomodoros)
			j += task.NCompleted()
		}
	})
	if i == 0 || j == 0 {
		return 0
	}
	return int((j / i) * 100)
}

func (p Project) Info() string {
	buf := bytes.NewBuffer(nil)
	pc := PercentComplete(p)
	fmt.Fprintf(buf, "[P%d]", p.ID)
	fmt.Fprintf(buf, "[")
	if pc == 100 {
		color.New(color.FgHiGreen).Fprintf(buf, "%d%%", pc)
	} else if pc > 50 && pc < 100 {
		color.New(color.FgHiYellow).Fprintf(buf, "%d%%", pc)
	} else {
		color.New(color.FgHiMagenta).Fprintf(buf, "%d%%", pc)
	}
	fmt.Fprintf(buf, "]")
	fmt.Fprintf(buf, "[%s]", truncDuration(Duration(p).String()))
	if p.Title != "" {
		fmt.Fprintf(buf, " %s", p.Title)
	}
	return buf.String()
}

const (
	emptySpace   = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

type Tree Project

func (t Tree) MaxDepth() int {
	depth := len(t.Children)
	for _, child := range t.Children {
		depth += Tree(*child).MaxDepth()
	}
	return depth
}

func (t Tree) Write(w io.Writer, depth int, last bool) {
	if depth == 0 { // root
		// fmt.Fprintf(w, ".\n")
		fmt.Fprintf(w, "%s\n", Project(t).Info())
	}
	spaces := depth
	// task list
	for i, task := range t.Tasks {
		for j := 0; j < spaces; j++ {
			if j == 0 && !last {
				fmt.Fprintf(w, continueItem)
				continue
			}
			fmt.Fprintf(w, emptySpace)
		}
		if i+1 == len(t.Tasks) && t.MaxDepth() == 0 {
			fmt.Fprintf(w, lastItem)
		} else {
			fmt.Fprintf(w, middleItem)
		}
		fmt.Fprintf(w, "%s\n", task.Info())
		// pomodoro list
		if len(task.Pomodoros) > 0 {
			for j := 0; j < spaces+1; j++ {
				if j == 0 && !last {
					fmt.Fprintf(w, continueItem)
					continue
				}
				if j == spaces && (i != len(t.Tasks)-1 || t.MaxDepth() > 0) {
					fmt.Fprintf(w, continueItem)
				}
				fmt.Fprintf(w, emptySpace)
			}
			fmt.Fprintf(w, lastItem)
			for _, p := range task.Pomodoros {
				// fmt.Fprintf(w, "[PM%d]", k)
				fmt.Fprintf(w, "%s", p.Info(task.Duration))
			}
			fmt.Fprintf(w, "\n")
		}
	}

	// sub projects
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
		fmt.Fprintf(w, "%s\n", child.Info())
		Tree(*child).Write(w, depth+1, n+1 == len(t.Children) && depth == 0)
	}
}
