package main

import (
	"os"
	// "sort"
	"encoding/json"
	"fmt"

	cli "github.com/jawher/mow.cli"
)

func get(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.LongDesc = `
Examples:

# output all tasks across all projects
pomo get
# output all tasks across all projects as a tree
pomo get --tree
# output tasks from a particular project filtered by name
pomo get --project "my project name"
# output tasks from a particular project filtered by tag
pomo get --project key=value
# output tasks matching the filter for a particular project
pomo get --project "another project" --task "some task name"
# complex task matching for a particular project
pomo get --project "another project" --task "key=value" --task "fuu=bar"
        `
		var (
			asTree = cmd.BoolOpt("t tree", false, "write projects and tasks as a tree")
			// otherArgs = cmd.StringsArg("ARG", []string{}, "filters")
			projectFilterArgs = cmd.StringsOpt("p project", []string{}, "project filters")
			taskFilterArgs    = cmd.StringsOpt("t task", []string{}, "task filters")
			// taskFilters    = cmd.StringsOpt("t task", []string{}, "task filters")
		)
		cmd.Action = func() {
			root := &Project{
				ID: int64(0),
			}
			projects := []*Project{root}
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				return s.ReadProject(root)
			}))
			projects = FilterProjects(projects, ProjectFiltersFromStrings(*projectFilterArgs)...)
			for _, project := range projects {
				ForEachMutate(project, func(p *Project) {
					p.Tasks = FilterTasks(p.Tasks, TaskFiltersFromStrings(*taskFilterArgs)...)
				})
			}
			projects = FilterProjects(projects, ProjectFilterSomeTasks())
			if config.JSON {
				maybe(json.NewEncoder(os.Stdout).Encode(projects))
				return
			}
			for _, project := range projects {
				if *asTree {
					Tree(*project).Write(os.Stdout, 0, Tree(*project).MaxDepth() == 0)
				} else {
					ForEach(*project, func(p Project) {
						fmt.Fprintln(os.Stdout, p.Info())
						for _, task := range p.Tasks {
							fmt.Fprintf(os.Stdout, " %s\n", task.Info())
						}
					})
				}
			}
		}
	}
}
