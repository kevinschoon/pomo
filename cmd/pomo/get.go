package main

import (
	"os"
	// "sort"
	"encoding/json"
	"fmt"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/filter"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tree"
)

func get(cfg *config.Config) func(*cli.Cmd) {
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
			asTree = cmd.BoolOpt("tree", false, "write projects and tasks as a tree")
			// otherArgs = cmd.StringsArg("ARG", []string{}, "filters")
			projectFilterArgs = cmd.StringsOpt("p project", []string{}, "project filters")
			taskFilterArgs    = cmd.StringsOpt("t task", []string{}, "task filters")
			// taskFilters    = cmd.StringsOpt("t task", []string{}, "task filters")
		)
		cmd.Action = func() {
			root := &pomo.Project{
				ID: int64(0),
			}
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				return db.ReadProject(root)
			}))

			filters := filter.Filters{
				ProjectFilters: filter.ProjectFiltersFromStrings(*projectFilterArgs),
				TaskFilters:    filter.TaskFiltersFromStrings(*taskFilterArgs),
			}

			root = filter.Reduce(root, filters)

			if cfg.JSON {
				maybe(json.NewEncoder(os.Stdout).Encode(root.Children))
				return
			}
			for _, project := range root.Children {
				if *asTree {
					tree.Tree(*project).Write(os.Stdout, 0, tree.Tree(*project).MaxDepth() == 0)
				} else {
					pomo.ForEach(*project, func(p pomo.Project) {
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
