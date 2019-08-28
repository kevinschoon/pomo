package main

import (
	"fmt"
	"os"
	// "sort"
	"encoding/json"

	cli "github.com/jawher/mow.cli"
)

type Kind string

const (
	PROJECT  = "projects"
	TASK     = "tasks"
	POMODORO = "pomodoros"
)

func getProject(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID"
		var (
			id = cmd.IntArg("ID", 0, "project identifier")
		)
		cmd.Action = func() {
			project := &Project{
				ID: int64(*id),
			}
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				return s.ReadProject(project)
			}))
			if config.JSON {
				maybe(json.NewEncoder(os.Stdout).Encode(project))
				return
			}
			Tree(*project).Write(os.Stdout, 0, Tree(*project).MaxDepth() == 0)
		}
	}
}

func getTask(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID"
		var (
			id = cmd.IntArg("ID", 0, "project identifier")
		)
		cmd.Action = func() {

			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			var tasks []*Task
			maybe(store.With(func(s Store) error {
				result, err := s.ReadTasks(int64(*id))
				if err != nil {
					return err
				}
				tasks = result
				return nil
			}))
			for _, task := range tasks {
				fmt.Printf("%s\n", task.Info())
			}
		}
	}
}

func get(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Command("project", "get a project", getProject(config))
		cmd.Command("task", "get a task", getTask(config))
		cmd.Action = func() {
			project := &Project{
				ID: int64(0),
			}
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				return s.ReadProject(project)
			}))
			if config.JSON {
				maybe(json.NewEncoder(os.Stdout).Encode(project))
				return
			}
			Tree(*project).Write(os.Stdout, 0, Tree(*project).MaxDepth() == 0)
		}
	}
}
