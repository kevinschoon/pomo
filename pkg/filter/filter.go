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

func (f Filters) Find(tasks []*pomo.Task) []*pomo.Task {
	if f.Empty() {
		return tasks
	}
	var filtered []*pomo.Task
	for _, task := range tasks {
		pomo.ForEachMutate(task, func(subTask *pomo.Task) {
			if f.MatchAll(*subTask) {
				filtered = append(filtered, subTask)
			}
		})
	}
	return filtered
}

func (f Filters) FindOne(tasks []*pomo.Task) (*pomo.Task, error) {
	var (
		result *pomo.Task
		err    error
	)
	for _, task := range tasks {
		pomo.ForEach(*task, func(subTask pomo.Task) {
			if f.MatchAll(subTask) {
				if result != nil {
					err = ErrTooManyResults
				}
			}
		})
	}
	if result == nil {
		return nil, ErrNoResults
	}
	return result, err
}
