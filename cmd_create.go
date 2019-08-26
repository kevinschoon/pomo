package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"
)

func createProject(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			path   = cmd.StringOpt("path", "", "path to a project file")
			title  = cmd.StringOpt("t title", "", "project title")
			parent = cmd.IntOpt("p parent", 0, "parent project id")
		)
		cmd.Action = func() {
			project := &Project{}
			if *path != "" {
				if *path == "-" {
					maybe(json.NewDecoder(os.Stdin).Decode(project))
				} else {
					raw, err := ioutil.ReadFile(*path)
					maybe(err)
					maybe(json.Unmarshal(raw, project))
				}
			}
			if *title != "" {
				project.Title = *title
			}
			if *parent != 0 {
				project.ParentID = int64(*parent)
			}
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				return s.CreateProject(project)
			}))
		}
	}
}

func createTask(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] MESSAGE"
		var (
			projectId = cmd.IntOpt("project", 0, "project id")
			duration  = cmd.StringOpt("d duration", "25m", "duration of each stent")
			pomodoros = cmd.IntOpt("p pomodoros", 4, "number of pomodoros")
			message   = cmd.StringArg("MESSAGE", "", "descriptive name of the given task")
			tags      = cmd.StringsOpt("t tag", []string{}, "tags associated with this task")
		)
		cmd.Action = func() {
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			task := &Task{
				ProjectID: int64(*projectId),
				Message:   *message,
				Duration:  parsed,
				Pomodoros: NewPomodoros(*pomodoros),
				Tags:      *tags,
			}
			maybe(store.With(func(s Store) error {
				return store.CreateTask(task)
			}))
			fmt.Printf("%d", task.ID)
		}
	}
}

func create(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("task", "create a new task", createTask(config))
		cmd.Command("project", "create a new project", createProject(config))
	}
}
