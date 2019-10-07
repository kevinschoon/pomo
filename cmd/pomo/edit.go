package main

import (
	"fmt"
	"time"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
)

func edit(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID"
		cmd.LongDesc = `
Update an existing task
        `
		var (
			taskID       = cmd.IntArg("ID", -1, "task identifier")
			parentID     = cmd.IntOpt("parent", -1, "parent id")
			duration     = cmd.StringOpt("duration", "", "pomodoro duration")
			addPomodoros = cmd.IntOpt("a add", 0, "add n pomodoros")
			delPomodoro  = cmd.IntOpt("d del", -1, "delete pomodoro")
			message      = cmd.StringOpt("m message", "", "message")
			truncate     = cmd.BoolOpt("truncate", false, "truncate the task to it's current runtime")
			done         = cmd.BoolOpt("done", false, "mark the task as completed")
			kvs          = cmd.StringsOpt("t tag", []string{}, "project tags")
		)
		cmd.Action = func() {
			tgs, err := tags.FromKVs(*kvs)
			maybe(err)
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			if *addPomodoros > 0 && *delPomodoro != -1 {
				maybe(fmt.Errorf("cannot add and delete pomodoros in one operation"))
			}
			maybe(db.With(func(db store.Store) error {
				err := db.Snapshot()
				if err != nil {
					return err
				}
				task := &pomo.Task{
					ID: int64(*taskID),
				}
				err = db.ReadTask(task)
				if err != nil {
					return err
				}
				if *message != "" {
					task.Message = *message
				}
				if *parentID != -1 {
					task.ParentID = int64(*parentID)
				}
				if *duration != "" {
					parsed, err := time.ParseDuration(*duration)
					maybe(err)
					task.Duration = parsed
				}
				tags.Merge(task.Tags, tgs)
				if *addPomodoros > 0 {
					for _, pomodoro := range pomo.NewPomodoros(*addPomodoros) {
						pomodoro.TaskID = task.ID
						err = db.CreatePomodoro(pomodoro)
						if err != nil {
							return err
						}
					}
				}
				if *delPomodoro != -1 {
					if *delPomodoro+1 > len(task.Pomodoros) {
						return fmt.Errorf("no pomodoro at index %d", *delPomodoro)
					}
					targetID := task.Pomodoros[*delPomodoro].ID
					err = db.DeletePomodoros(task.ID, targetID)
					if err != nil {
						return err
					}
				}
				if *truncate {
					task.Truncate()
					for _, pomodoro := range task.Pomodoros {
						err = db.UpdatePomodoro(pomodoro)
						if err != nil {
							return err
						}
					}
				}
				if *done {
					task.Fill()
					for _, pomodoro := range task.Pomodoros {
						err = db.UpdatePomodoro(pomodoro)
						if err != nil {
							return err
						}
					}
				}
				return db.UpdateTask(task)
			}))

		}
	}
}
