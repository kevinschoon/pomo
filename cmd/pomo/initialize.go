package main

import (
	cli "github.com/jawher/mow.cli"

	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
)

func initialize(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			maybe(db.Init())
		}
	}
}
