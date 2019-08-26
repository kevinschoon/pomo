package main

import (
	cli "github.com/jawher/mow.cli"
)

func initialize(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			db, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(initDB(db))
		}
	}
}
