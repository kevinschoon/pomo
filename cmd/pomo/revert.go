package main

import (
	"encoding/json"
	"os"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
)

func history(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [ID]"
		var (
			id      = cmd.IntArg("ID", 0, "snapshot id")
			dump    = cmd.BoolOpt("d dump", false, "dump the previous state")
			restore = cmd.BoolOpt("r restore", false, "restore the snapshot")
		)
		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				last := pomo.NewTask()
				err = db.Revert(*id, last)
				if err != nil {
					return err
				}
				if *dump {
					json.NewEncoder(os.Stdout).Encode(last)
				}
				if *restore {
					err = db.Reset()
					if err != nil {
						return err
					}
					for _, task := range last.Tasks {
						err := db.WriteTask(task)
						if err != nil {
							return err
						}
					}
				}
				return err
			}))

		}
	}
}
