package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
)

func createProject(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [TITLE]"
		cmd.LongDesc = `
Create a new project by specifying a title or path to a JSON file
containing a configuration. If - is specified as an argument to --path
configuration will be read from stdin.
        `
		var (
			title  = cmd.StringArg("TITLE", "", "project title")
			path   = cmd.StringOpt("path", "", "path to a project file")
			parent = cmd.IntOpt("p parent", 0, "parent project id")
			kvs    = cmd.StringsOpt("t tag", []string{}, "project tags")
		)
		cmd.Action = func() {
			project := &pomo.Project{}
			tgs, err := tags.FromKVs(*kvs)
			project.Tags = tgs
			maybe(err)
			if *path != "" {
				if *path == "-" {
					maybe(json.NewDecoder(os.Stdin).Decode(project))
				} else {
					raw, err := ioutil.ReadFile(*path)
					maybe(err)
					maybe(json.Unmarshal(raw, project))
				}
			} else {
				if *title == "" {
					maybe(fmt.Errorf("need to specify a title or project file"))
				}
			}
			if *title != "" {
				project.Title = *title
			}
			if *parent != 0 {
				project.ParentID = int64(*parent)
			}
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				return db.CreateProject(project)
			}))
		}
	}
}

func createTask(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] MESSAGE"
		var (
			projectId = cmd.IntOpt("project", 0, "project id")
			duration  = cmd.StringOpt("d duration", "25m", "duration of each stent")
			pomodoros = cmd.IntOpt("p pomodoros", 4, "number of pomodoros")
			message   = cmd.StringArg("MESSAGE", "", "descriptive name of the given task")
			kvs       = cmd.StringsOpt("t tag", []string{}, "tags associated with this task")
		)
		cmd.Action = func() {
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			tgs, err := tags.FromKVs(*kvs)
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			maybe(err)
			task := &pomo.Task{
				ProjectID: int64(*projectId),
				Message:   *message,
				Duration:  parsed,
				Pomodoros: pomo.NewPomodoros(*pomodoros),
				Tags:      tgs,
			}
			maybe(db.With(func(db store.Store) error {
				return db.CreateTask(task)
			}))
			fmt.Printf("%d", task.ID)
		}
	}
}

func create(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("t task", "create a new task", createTask(cfg))
		cmd.Command("p project", "create a new project", createProject(cfg))
	}
}
