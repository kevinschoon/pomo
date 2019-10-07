package main

import (
	"fmt"
	"strconv"
	"time"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/harness"
	"github.com/kevinschoon/pomo/pkg/notify"
	"github.com/kevinschoon/pomo/pkg/runner"
	"github.com/kevinschoon/pomo/pkg/runner/server"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
	"github.com/kevinschoon/pomo/pkg/ui"
)

func start(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [TASK]"
		var (
			create          = cmd.BoolOpt("c create", false, "create a new task before starting")
			taskDescription = cmd.StringArg("TASK", "", "task id or message (with --create)")
			parentID        = cmd.IntOpt("parent", 0, "parent task id")
			pomodoros       = cmd.IntOpt("p pomodoros", cfg.DefaultPomodoros, "number of pomodoros")
			durationStr     = cmd.StringOpt("d duration", cfg.DefaultDuration.String(), "task duration")
			kvs             = cmd.StringsOpt("t tag", []string{}, "task tags")
		)
		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			task := &pomo.Task{}
			if *create {
				if *taskDescription == "" {
					maybe(fmt.Errorf("need to provide a task description"))
				}
				task.Message = *taskDescription
				task.ParentID = int64(*parentID)
				duration, err := time.ParseDuration(*durationStr)
				maybe(err)
				task.Duration = duration
				task.Pomodoros = pomo.NewPomodoros(*pomodoros)
				tgs, err := tags.FromKVs(*kvs)
				maybe(err)
				task.Tags = tgs
				maybe(db.With(func(db store.Store) error {
					return db.CreateTask(task)
				}))
			} else {
				taskID, err := strconv.ParseUint(*taskDescription, 0, 64)
				if err != nil {
					maybe(fmt.Errorf("cannot parse taskID: %s", err.Error()))
				}
				task.ID = int64(taskID)
				maybe(db.With(func(db store.Store) error {
					return db.ReadTask(task)
				}))
			}
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
