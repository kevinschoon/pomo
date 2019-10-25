package tree

import (
	"bytes"
	"fmt"
	"io"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config/color"
)

const (
	emptySpace   = "    "
	middleItem   = "├──"
	continueItem = "│   "
	lastItem     = "└──"
)

// Tree prints a hierarchy of tasks
// in a tree.
type Tree struct {
	pomo.Task
	Colors        color.Colors
	ShowPomodoros bool
	TaskTemplater func(pomo.Task) string
}

func New(task pomo.Task, pomodoros bool) Tree {
	return Tree{
		Task:          task,
		ShowPomodoros: pomodoros,
	}
}

func (t Tree) next(value bool, depth []bool) (result []bool) {
	return append(append(result, depth...), value)
}

func (t Tree) fill(w io.Writer, depth []bool) {
	for i := 0; i < len(depth); i++ {
		if depth[i] {
			t.Colors.Primary.Fprint(w, continueItem)
		} else {
			fmt.Fprint(w, emptySpace)
		}
	}
}

// Write writes the Tree representation of a Task hierarchy
// to the io.Writer
func (t Tree) Write(w io.Writer, depth []bool) {
	rootComplete := pomo.Complete(t.Task)
	if depth == nil { // root
		if rootComplete {
			t.Colors.Tertiary.Fprintf(w, "%s\n", t.TaskTemplater(t.Task))
		} else {
			t.Colors.Primary.Fprintf(w, "%s\n", t.TaskTemplater(t.Task))
		}
	}
	for n, task := range t.Tasks {
		taskComplete := pomo.Complete(*task)
		last := n+1 == len(t.Tasks)
		t.fill(w, depth)
		var item string
		if last {
			item = lastItem
			// fmt.Fprint(w, lastItem)
		} else {
			item = middleItem
			//fmt.Fprint(w, middleItem)
		}
		if rootComplete {
			t.Colors.Tertiary.Fprintf(w, "%s%s\n", item, t.TaskTemplater(*task))
		} else {
			t.Colors.Primary.Fprintf(w, item)
			if taskComplete {
				t.Colors.Tertiary.Fprintf(w, "%s\n", t.TaskTemplater(*task))
			} else {
				t.Colors.Primary.Fprintf(w, "%s\n", t.TaskTemplater(*task))
			}
		}
		next := Tree{
			Task:          *task,
			Colors:        t.Colors,
			ShowPomodoros: t.ShowPomodoros,
			TaskTemplater: t.TaskTemplater,
		}
		next.Write(w, t.next(len(task.Tasks) > 0 && !last, depth))
	}
}

func (t Tree) String() string {
	buf := bytes.NewBuffer(nil)
	t.Write(buf, nil)
	return buf.String()
}
