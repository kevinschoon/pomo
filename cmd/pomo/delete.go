package main

import (
	cli "github.com/jawher/mow.cli"

	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
)

func deleteTask(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID"
		var (
			taskID = cmd.IntArg("ID", -1, "task to delete")
		)

		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				err := db.Snapshot()
				if err != nil {
					return err
				}
				_, err = db.ReadTask(int64(*taskID))
				if err != nil {
					return err
				}
				return db.DeleteTask(int64(*taskID))
			}))
		}
	}
}
