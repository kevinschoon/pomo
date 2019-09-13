package main

import (
	"encoding/json"
	"os"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
)

func importTasks(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] PATH"
		var (
			path = cmd.StringArg("PATH", "", "path to import data, - for stdin")
		)

		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			root := &pomo.Task{}
			if *path == "-" {
				maybe(json.NewDecoder(os.Stdin).Decode(root))
			} else {
				fp, err := os.Open(*path)
				maybe(err)
				maybe(json.NewDecoder(fp).Decode(root))
			}

			maybe(db.With(func(db store.Store) error {
				var err error
				pomo.ForEachMutate(root, func(task *pomo.Task) {
					if err != nil {
						return
					}
					if task.ID == int64(0) {
						return
					}
					e := db.CreateTask(task)
					if e != nil {
						err = e
					}
					for _, subTask := range task.Tasks {
						subTask.ParentID = task.ID
					}
				})
				return err
			}))
		}
	}
}
