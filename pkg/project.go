package pomo

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/kevinschoon/pomo/pkg/internal/format"
	"github.com/kevinschoon/pomo/pkg/tags"
)

// Project is a logical container for tasks
type Project struct {
	ID       int64      `json:"id"`
	ParentID int64      `json:"parent_id"`
	Title    string     `json:"title"`
	Children []*Project `json:"children"`
	Tasks    []*Task    `json:"tasks"`
	Tags     *tags.Tags `json:"tags"`
}

func NewProject() *Project {
	return &Project{
		Tags: tags.New(),
	}
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
	fmt.Fprintf(buf, "[%s]", format.TruncDuration(Duration(p).String()))
	for _, key := range p.Tags.Keys() {
		if p.Tags.Get(key) == "" {
			fmt.Fprintf(buf, "[%s]", key)
		} else {
			fmt.Fprintf(buf, "[%s=%s]", key, p.Tags.Get(key))
		}
	}
	if p.Title != "" {
		fmt.Fprintf(buf, " %s", p.Title)
	}
	return buf.String()
}

func (p Project) FlattenTasks() []*Task {
	var tasks []*Task
	for _, task := range p.Tasks {
		tasks = append(tasks, task)
	}
	for _, child := range p.Children {
		for _, task := range child.FlattenTasks() {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

type ProjectFn func(Project)

func ForEach(p Project, fn ProjectFn) {
	fn(p)
	for _, child := range p.Children {
		ForEach(*child, fn)
	}
}

func ForEachMutate(p *Project, fn func(*Project)) {
	fn(p)
	for _, child := range p.Children {
		ForEachMutate(child, fn)
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
