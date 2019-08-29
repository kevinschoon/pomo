package main

import (
	"fmt"
	"time"

	cli "github.com/jawher/mow.cli"
)

func editProject(config *Config) func(*cli.Cmd) {
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
			project := &Project{
				ID: int64(*projectID),
			}
			tags, err := NewTagsFromKVs(*kvs)
			maybe(err)
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				err := s.ReadProject(project)
				if err != nil {
					return err
				}
				if *parentID != -1 {
					project.ParentID = int64(*parentID)
				}
				if *title != "" {
					project.Title = *title
				}
				MergeTags(project.Tags, tags)
				return s.UpdateProject(project)
			}))
		}
	}
}

func editTask(config *Config) func(*cli.Cmd) {
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
			tags, err := NewTagsFromKVs(*kvs)
			maybe(err)
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			if *addPomodoros > 0 && *delPomodoro != -1 {
				maybe(fmt.Errorf("cannot add and delete pomodoros in one operation"))
			}
			maybe(store.With(func(s Store) error {
				task := &Task{
					ID: int64(*taskID),
				}
				err := s.ReadTask(task)
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
				MergeTags(task.Tags, tags)
				if *addPomodoros > 0 {
					for _, pomodoro := range NewPomodoros(*addPomodoros) {
						pomodoro.TaskID = task.ID
						err = s.CreatePomodoro(pomodoro)
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
					err = s.DeletePomodoros(task.ID, targetID)
					if err != nil {
						return err
					}
				}
				return store.UpdateTask(task)
			}))

		}
	}
}

func edit(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("p project", "edit an existing project", editProject(config))
		cmd.Command("t task", "edit an existing task", editTask(config))
	}
}
