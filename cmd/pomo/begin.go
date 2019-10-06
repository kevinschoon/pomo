package main

import (
	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/harness"
	"github.com/kevinschoon/pomo/pkg/notify"
	"github.com/kevinschoon/pomo/pkg/runner"
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
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			task := &pomo.Task{
				ID: int64(*taskId),
			}
			maybe(db.With(func(db store.Store) error {
				err := db.Snapshot()
				if err != nil {
					return err
				}
				return db.ReadTask(task)
			}))
			notifier := notify.NewXNotifier(cfg.IconPath)
			statusCh := make(chan runner.Status, 20)
			socketServer := server.NewSocketServer(cfg.SocketPath)
			taskRunner := runner.NewTaskRunner(task,
				socketServer.SetStatus,
				runner.StatusTicker(statusCh),
				runner.StatusUpdater(task, db),
				notify.StatusFunc(notifier),
			)
			termUI := ui.New(taskRunner.Toggle, taskRunner.Suspend, statusCh)
			maybe(harness.Harness{
				UI:     termUI,
				Server: socketServer,
				Runner: taskRunner,
			}.Launch())
		}
	}
}
