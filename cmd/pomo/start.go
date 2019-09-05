package main

import (
	"time"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/harness"
	"github.com/kevinschoon/pomo/pkg/runner"
	"github.com/kevinschoon/pomo/pkg/runner/server"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
	"github.com/kevinschoon/pomo/pkg/ui"
)

func start(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] MESSAGE"
		var (
			duration  = cmd.StringOpt("d duration", "25m", "duration of each stent")
			pomodoros = cmd.IntOpt("p pomodoros", 4, "number of pomodoros")
			message   = cmd.StringArg("MESSAGE", "", "descriptive name of the given task")
			kvs       = cmd.StringsOpt("t tag", []string{}, "tags associated with this task")
		)
		cmd.Action = func() {
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			tgs, err := tags.FromKVs(*kvs)
			maybe(err)
			task := &pomo.Task{
				Message:   *message,
				Tags:      tgs,
				Pomodoros: pomo.NewPomodoros(*pomodoros),
				Duration:  parsed,
			}
			statusCh := make(chan runner.Status, 20)
			socketServer := server.NewSocketServer(cfg.SocketPath)
			taskRunner := runner.NewTaskRunner(task, runner.JoinStatusFuncs(
				socketServer.SetStatus,
				runner.StatusTicker(statusCh),
				runner.StatusUpdater(task, db),
			))
			termUI := ui.New(taskRunner.Toggle, taskRunner.Suspend, statusCh)
			maybe(harness.Harness{
				UI:     termUI,
				Server: socketServer,
				Runner: taskRunner,
			}.Launch())
		}
	}

}
