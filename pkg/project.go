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

func (p Project) Duration() time.Duration {
	var duration time.Duration
	for _, task := range p.Tasks {
		duration += task.Duration * time.Duration(int64(len(task.Pomodoros)))
	}
	for _, child := range p.Children {
		duration += child.Duration()
	}
	return duration
}

func (p Project) TimeRunning() time.Duration {
	var running time.Duration
	for _, task := range p.Tasks {
		running += task.TimeRunning()
	}
	for _, child := range p.Children {
		running += child.TimeRunning()
	}
	return running
}

func (p Project) PercentComplete() float64 {
	duration := p.Duration()
	timeRunning := p.TimeRunning()
	return (float64(timeRunning) / float64(duration)) * 100
}

func (p Project) Info() string {
	buf := bytes.NewBuffer(nil)
	pc := int(p.PercentComplete())
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
	fmt.Fprintf(buf, "[%s]", format.TruncDuration(p.Duration().String()))
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
