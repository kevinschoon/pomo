package functional

import (
	"errors"
	"strings"

	pomo "github.com/kevinschoon/pomo/pkg"
)

var (
	// ErrTooManyResults indicates more than one result was found
	ErrTooManyResults = errors.New("too many results")
	// ErrNoResults indicates no results were found
	ErrNoResults = errors.New("no results")
)

// Filter is a function for matching a Task
type Filter func(pomo.Task) bool

// FiltersFromStrings creates filters based on
// an array of strings, likely from the CLI.
func FiltersFromStrings(args []string) []Filter {
	var filters []Filter
	for _, arg := range args {
		split := strings.Split(arg, "=")
		if len(split) == 1 {
			filters = append(filters, FilterByAny(FilterByTag(split[0], ""), FilterByName(split[0])))
		} else if len(split) == 2 {
			filters = append(filters, FilterByTag(split[0], split[1]))
		}
	}
	return filters
}

// FilterByName filters a task by it's name
func FilterByName(name string) Filter {
	return func(t pomo.Task) bool {
		return strings.Contains(t.Message, name)
	}
}

// FilterByTag filters a task by a tag Key and Value
func FilterByTag(key, value string) Filter {
	return func(t pomo.Task) bool {
		return (t.Tags.HasTag(key) && t.Tags.Get(key) == value)
	}
}

// FilterByID filters a task by it's ID number
func FilterByID(id int64) Filter {
	return func(t pomo.Task) bool {
		return t.ID == id
	}
}

// FilterByAny filters by any matching filter
func FilterByAny(filters ...Filter) Filter {
	return func(t pomo.Task) bool {
		for _, filter := range filters {
			if filter(t) {
				return true
			}
		}
		return false
	}
}

// MatchAny returns true if any filters match the task
func MatchAny(t pomo.Task, filters ...Filter) bool {
	for _, filter := range filters {
		if filter(t) {
			return true
		}
	}
	return false
}

// MatchAll returns true if all filters match the task
func MatchAll(t pomo.Task, filters ...Filter) bool {
	for _, filter := range filters {
		if !filter(t) {
			return false
		}
	}
	return true
}

// FindMany reduces the tasks to matching results
func FindMany(tasks []*pomo.Task, filters ...Filter) []*pomo.Task {
	if len(filters) == 0 {
		return tasks
	}
	var filtered []*pomo.Task
	for _, task := range tasks {
		ForEachMutate(task, func(subTask *pomo.Task) {
			if MatchAll(*subTask, filters...) {
				filtered = append(filtered, subTask)
			}
		})
	}
	return filtered
}

// FineOne finds a single task returning an error if none or more
// than one task are found
func FindOne(tasks []*pomo.Task, filters ...Filter) (*pomo.Task, error) {
	var (
		result *pomo.Task
		err    error
	)
	for _, task := range tasks {
		ForEach(*task, func(subTask pomo.Task) {
			if MatchAll(subTask, filters...) {
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
