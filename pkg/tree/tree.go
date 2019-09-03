package tree

import (
	"fmt"
	"io"

	pomo "github.com/kevinschoon/pomo/pkg"
)

const (
	emptySpace   = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

type Tree pomo.Project

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
		fmt.Fprintf(w, "%s\n", pomo.Project(t).Info())
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
