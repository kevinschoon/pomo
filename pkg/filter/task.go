package filter

import (
	"strings"

	pomo "github.com/kevinschoon/pomo/pkg"
)

type TaskFilter func(pomo.Task) bool

func TaskFiltersFromStrings(args []string) []TaskFilter {
	var filters []TaskFilter
	for _, arg := range args {
		split := strings.Split(arg, "=")
		if len(split) == 1 {
			filters = append(filters, TaskFilterByAny(TaskFilterByTag(split[0], ""), TaskFilterByName(split[0])))
		} else if len(split) == 2 {
			filters = append(filters, TaskFilterByTag(split[0], split[1]))
		}
	}
	return filters
}

func TaskFilterByName(name string) TaskFilter {
	return func(t pomo.Task) bool {
		return strings.Contains(t.Message, name)
	}
}

func TaskFilterByTag(key, value string) TaskFilter {
	return func(t pomo.Task) bool {
		return (t.Tags.HasTag(key) && t.Tags.Get(key) == value)
	}
}

func TaskFilterByID(id int64) TaskFilter {
	return func(t pomo.Task) bool {
		return t.ID == id
	}
}

func TaskFilterByAny(filters ...TaskFilter) TaskFilter {
	return func(t pomo.Task) bool {
		for _, filter := range filters {
			if filter(t) {
				return true
			}
		}
		return false
	}
}
