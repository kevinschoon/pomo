package filter

import (
	"errors"

	pomo "github.com/kevinschoon/pomo/pkg"
)

var (
	ErrTooManyResults = errors.New("too many results")
	ErrNoResults      = errors.New("no results")
)

type Filters []TaskFilter

func (f Filters) Empty() bool {
	return len(f) == 0
}

func (f Filters) MatchAny(item pomo.Task) bool {
	for _, filter := range f {
		if filter(item) {
			return true
		}
	}
	return false
}

func (f Filters) MatchAll(item pomo.Task) bool {
	for _, filter := range f {
		if !filter(item) {
			return false
		}
	}
	return true
}

func (f Filters) Reduce(root *pomo.Task) *pomo.Task {
	return root
}

func (f Filters) FindOne(root *pomo.Task) (*pomo.Task, error) {
	return nil, ErrNoResults
}

/*

func Reduce(root *pomo.Task, filters Filters) *pomo.Task {
	if filters.Empty() {
		return root
	}
	// map of all children that contain matches
	m := map[int64]bool{}
	pomo.ForEach(*root, func(task pomo.Task) {
		m[task.ID] = false
	})
	pomo.ForEach(*root, func(task pomo.Task) {
		pomo.ForEach(task, func(child pomo.Task) {
			if filters.MatchAny(child) {
				m[task.ID] = true
			}
		})
		// all decendents of a match are included
		if m[task.ID] {
			pomo.ForEach(task, func(child pomo.Task) {
				m[child.ID] = true
			})
		}
	})
	pomo.ForEachMutate(root, func(task *pomo.Task) {
		var (
			children []*pomo.Task
		)
		for _, child := range task.Tasks {
			if m[child.ID] == true {
				children = append(children, child)
			}
		}
		task.Tasks = children
	})
	return root
}

func FindOne(root pomo.Task, filters Filters) (pomo.Task, error) {
		pomo.ForEach(root, func(project pomo.Project) {
			if result.Error() != nil {
				return
			}
			if filters.MatchAll(Item{Project: &project}) {
				if result.project != nil || result.task != nil {
					result.err = ErrTooManyResults
				}
				result.project = &project
			}
			for _, task := range project.Tasks {
				if filters.MatchAll(Item{Task: task}) {
					if result.project != nil || result.task != nil {
						result.err = ErrTooManyResults
					}
					result.task = task
				}
			}
		})
		if result.Project() == nil && result.Task() == nil {
			result.err = ErrNoResults
		}
		return *result
	return root, nil // TODO
}

*/
