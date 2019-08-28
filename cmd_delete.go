package main

import (
	cli "github.com/jawher/mow.cli"
)

func deleteTask(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] TASK_ID"
		var taskID = cmd.IntArg("TASK_ID", -1, "task to delete")
		cmd.Action = func() {
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				return s.DeleteTask(int64(*taskID))
			}))
		}
	}
}

func deleteProject(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] PROJECT_ID"
		var projectID = cmd.IntArg("PROJECT_ID", -1, "project to delete")
		cmd.Action = func() {
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				return s.DeleteProject(int64(*projectID))
			}))
		}
	}
}

func _delete(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("task", "delete a task", deleteTask(config))
		cmd.Command("project", "delete a project", deleteProject(config))
	}
}
