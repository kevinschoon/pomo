package filter

import (
	"strings"

	pomo "github.com/kevinschoon/pomo/pkg"
)

type ProjectFilter func(pomo.Project) bool

func ProjectFiltersFromStrings(args []string) []ProjectFilter {
	var filters []ProjectFilter
	for _, arg := range args {
		split := strings.Split(arg, "=")
		if len(split) == 1 {
			filters = append(filters, ProjectFilterByTag(split[0], ""))
			filters = append(filters, ProjectFilterByName(split[0]))
		} else if len(split) == 2 {
			filters = append(filters, ProjectFilterByTag(split[0], split[1]))
		}
	}
	return filters
}

func ProjectFilterByName(name string) ProjectFilter {
	return func(p pomo.Project) bool {
		return strings.Contains(p.Title, name)
	}
}

func ProjectFilterByTag(key, value string) ProjectFilter {
	return func(p pomo.Project) bool {
		return (p.Tags.HasTag(key) && p.Tags.Get(key) == value)
	}
}

func ProjectFitlerByID(id int64) ProjectFilter {
	return func(p pomo.Project) bool {
		return p.ID == id
	}
}
