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
        `
		var (
			asTree  = cmd.BoolOpt("tree", false, "write projects and tasks as a tree")
			filters = cmd.StringsOpt("f filter", []string{}, "filters")
		)
		cmd.Action = func() {
			root := &pomo.Task{
				ID: int64(0),
			}
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				return db.ReadTask(root)
			}))

			filters := filter.Filters(filter.TaskFiltersFromStrings(*filters))
			root = filters.Reduce(root)

			if cfg.JSON {
				maybe(json.NewEncoder(os.Stdout).Encode(root.Tasks))
				return
			}
			for _, task := range root.Tasks {
				if *asTree {
					tree.Tree(*task).Write(os.Stdout, 0, tree.Tree(*task).MaxDepth() == 0)
				} else {
					pomo.ForEach(*task, func(t pomo.Task) {
						fmt.Fprintln(os.Stdout, t.Info())
						for _, task := range t.Tasks {
							fmt.Fprintf(os.Stdout, " %s\n", task.Info())
						}
					})
				}
			}
		}
	}
}
