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

func FilterProjects(projects []*pomo.Project, filters ...ProjectFilter) []*pomo.Project {
	if len(filters) == 0 {
		return projects
	}
	var filtered []*pomo.Project
	match := func(p pomo.Project) bool {
		for _, filter := range filters {
			if filter(p) {
				return true
			}
		}
		return false
	}
	for _, project := range projects {
		if match(*project) {
			// direct match
			// project.Children = FilterProjects(project.Children, filters...)
			filtered = append(filtered, project)
		} else {
			for _, indirect := range FilterProjects(project.Children, filters...) {
				filtered = append(filtered, indirect)
			}
			/*
				indirect := FilterProjects(project.Children, filters...)
				if len(indirect) > 0 {
					project.Children = indirect
					filtered = append(filtered, project)
				}
			*/
		}
	}
	return filtered
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

func ProjectFilterSomeTasks() ProjectFilter {
	return func(p pomo.Project) bool {
		return len(p.Tasks) > 0
	}
}
