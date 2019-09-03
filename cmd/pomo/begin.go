package main

import (
	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/runner/server"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/ui"
)

func begin(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] TASK_ID"
		var (
			taskId = cmd.IntArg("TASK_ID", -1, "ID of Pomodoro to begin")
		)

		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			task := &pomo.Task{
				ID: int64(*taskId),
			}
			maybe(db.With(func(db store.Store) error {
				return db.ReadTask(task)
			}))
			server, err := server.NewSocketServer(task, db, cfg)
			maybe(err)
			shutdown := make(chan error)
			go func() {
				shutdown <- server.Serve()
			}()
			// runner.Start(task)
			// defer server.Stop()
			maybe(ui.Start(server))
			maybe(<-shutdown)
		}
	}
}
