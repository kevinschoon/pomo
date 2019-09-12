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
		cmd.Spec = "[OPTIONS] [FILTER...]"
		cmd.LongDesc = `
Examples:

# output all tasks across all projects
pomo get
# output all tasks across all projects as a tree
pomo get --tree
        `
		var (
			asTree        = cmd.BoolOpt("tree", true, "write projects and tasks as a tree")
			showPomodoros = cmd.BoolOpt("p pomodoros", true, "show status of each pomodoro")
			filters       = cmd.StringsArg("FILTER", []string{}, "filters")
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
			root.Tasks = filters.Find(root.Tasks)

			if cfg.JSON {
				maybe(json.NewEncoder(os.Stdout).Encode(root))
				return
			}
			if *asTree {
				tree.Tree{Task: *root, ShowPomodoros: *showPomodoros}.Write(os.Stdout, nil)
				// tree.Tree(*root).Write(os.Stdout, nil)
			} else {
				pomo.ForEach(*root, func(t pomo.Task) {
					fmt.Fprintln(os.Stdout, t.Info())
				})
			}
		}
	}
}
