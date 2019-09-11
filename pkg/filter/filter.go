package filter

import (
	"errors"

	pomo "github.com/kevinschoon/pomo/pkg"
)

var (
	ErrTooManyResults = errors.New("too many results")
	ErrNoResults      = errors.New("no results")
)

// pomo.Project | pomo.Task
type Item struct {
	Project *pomo.Project
	Task    *pomo.Task
}

// ProjectFilter | TaskFilter
type Filters struct {
	ProjectFilters []ProjectFilter
	TaskFilters    []TaskFilter
}

func (f Filters) Empty() bool {
	return len(f.ProjectFilters) == 0 && len(f.TaskFilters) == 0
}

func (f Filters) MatchAny(item Item) bool {
	if item.Project != nil {
		for _, projectFilter := range f.ProjectFilters {
			if projectFilter(*item.Project) {
				return true
			}
		}
	} else if item.Task != nil {
		for _, taskFilter := range f.TaskFilters {
			if taskFilter(*item.Task) {
				return true
			}
		}
	}
	return false
}

func (f Filters) MatchAll(item Item) bool {
	if item.Project != nil && len(f.ProjectFilters) > 0 {
		for _, projectFilter := range f.ProjectFilters {
			if projectFilter(*item.Project) {
				continue
			}
			return false
		}
		return true
	} else if item.Task != nil && len(f.TaskFilters) > 0 {
		for _, taskFilter := range f.TaskFilters {
			if taskFilter(*item.Task) {
				continue
			}
			return false
		}
		return true
	}

	return false
}

// pomo.Project | pomo.Task | (ErrTooManyResults | ErrNoResults | nil)
type Result struct {
	project *pomo.Project
	task    *pomo.Task
	err     error
}

func (r Result) Project() *pomo.Project {
	return r.project
}

func (r Result) Task() *pomo.Task {
	return r.task
}

func (r Result) Error() error {
	return r.err
}

func Reduce(root *pomo.Project, filters Filters) *pomo.Project {
	if filters.Empty() {
		return root
	}
	// map of all children that contain matches
	m := map[int64]bool{}
	pomo.ForEach(*root, func(project pomo.Project) {
		m[project.ID] = false
	})
	pomo.ForEach(*root, func(project pomo.Project) {
		pomo.ForEach(project, func(child pomo.Project) {
			for _, task := range project.Tasks {
				if filters.MatchAny(Item{Task: task}) {
					m[project.ID] = true
				}
			}
			if filters.MatchAny(Item{Project: &project}) {
				m[project.ID] = true
			}
		})
		// all decendents of a match are included
		if m[project.ID] {
			pomo.ForEach(project, func(child pomo.Project) {
				m[child.ID] = true
			})
		}
	})
	pomo.ForEachMutate(root, func(project *pomo.Project) {
		var (
			children []*pomo.Project
		)
		for _, child := range project.Children {
			if m[child.ID] == true {
				children = append(children, child)
			}
		}
		if len(filters.TaskFilters) > 0 {
			var tasks []*pomo.Task
			for _, task := range project.Tasks {
				if filters.MatchAny(Item{Task: task}) {
					tasks = append(tasks, task)
				}
			}
			project.Tasks = tasks
		}
		project.Children = children
	})
	return root
}

func FindOne(root pomo.Project, filters Filters) Result {
	result := &Result{}
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
}
