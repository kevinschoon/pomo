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

func editProject(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID"
		var (
			projectID = cmd.IntArg("ID", -1, "project identifier")
			parentID  = cmd.IntOpt("p parent", -1, "parent identifier")
			title     = cmd.StringOpt("t title", "", "title")
			kvs       = cmd.StringsOpt("t tag", []string{}, "project tags")
		)
		cmd.Action = func() {
			if *projectID == 0 {
				maybe(fmt.Errorf("root project may not be modified"))
			}
			project := &pomo.Project{
				ID: int64(*projectID),
			}
			tgs, err := tags.FromKVs(*kvs)
			maybe(err)
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				err := db.ReadProject(project)
				if err != nil {
					return err
				}
				if *parentID != -1 {
					project.ParentID = int64(*parentID)
				}
				if *title != "" {
					project.Title = *title
				}
				tags.Merge(project.Tags, tgs)
				return db.UpdateProject(project)
			}))
		}
	}
}

func editTask(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID"
		cmd.LongDesc = `
Update an existing task
        `
		var (
			taskID       = cmd.IntArg("ID", -1, "task identifier")
			projectID    = cmd.IntOpt("project-id", -1, "project identifier")
			duration     = cmd.StringOpt("duration", "", "pomodoro duration")
			addPomodoros = cmd.IntOpt("a add", 0, "add n pomodoros")
			delPomodoro  = cmd.IntOpt("d del", -1, "delete pomodoro")
			message      = cmd.StringOpt("m message", "", "message")
			kvs          = cmd.StringsOpt("t tag", []string{}, "project tags")
		)
		cmd.Action = func() {
			tgs, err := tags.FromKVs(*kvs)
			maybe(err)
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			if *addPomodoros > 0 && *delPomodoro != -1 {
				maybe(fmt.Errorf("cannot add and delete pomodoros in one operation"))
			}
			maybe(db.With(func(db store.Store) error {
				task := &pomo.Task{
					ID: int64(*taskID),
				}
				err := db.ReadTask(task)
				if err != nil {
					return err
				}
				if *message != "" {
					task.Message = *message
				}
				if *projectID != -1 {
					task.ProjectID = int64(*projectID)
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
				return db.UpdateTask(task)
			}))

		}
	}
}

func edit(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("p project", "edit an existing project", editProject(cfg))
		cmd.Command("t task", "edit an existing task", editTask(cfg))
	}
}
