package main

import (
	"strings"
)

func FilterTasks(tasks []*Task, filters ...TaskFilter) []*Task {
	if len(filters) == 0 {
		return tasks
	}
	var filtered []*Task
	for _, task := range tasks {
		for _, filter := range filters {
			if filter(*task) {
				filtered = append(filtered, task)
			}
		}
	}
	return filtered
}

type TaskFilter func(Task) bool

func TaskFiltersFromStrings(args []string) []TaskFilter {
	var filters []TaskFilter
	for _, arg := range args {
		split := strings.Split(arg, "=")
		if len(split) == 1 {
			filters = append(filters, TaskFilterByTag(split[0], ""))
			filters = append(filters, TaskFilterByName(split[0]))
		} else if len(split) == 2 {
			filters = append(filters, TaskFilterByTag(split[0], split[1]))
		}
	}
	return filters
}

func TaskFilterByName(name string) TaskFilter {
	return func(t Task) bool {
		return strings.Contains(t.Message, name)
	}
}

func TaskFilterByTag(key, value string) TaskFilter {
	return func(t Task) bool {
		return (t.Tags.HasTag(key) && t.Tags.Get(key) == value)
	}
}

func TaskFilterByID(id int64) TaskFilter {
	return func(t Task) bool {
		return t.ID == id
	}
}
