package main

import (
	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
)

func revert(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [ID]"
		var (
			id = cmd.IntArg("ID", 0, "snapshot id")
		)
		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				root := pomo.NewTask()
				err = db.ReadTask(root)
				if err != nil {
					return err
				}
				last := pomo.NewTask()
				err = db.Revert(*id, last)
				if err != nil {
					return err
				}
				err = db.Reset()
				if err != nil {
					return err
				}
				pomo.ForEachMutate(last, func(task *pomo.Task) {
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
