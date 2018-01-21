package main

import (
	"encoding/json"
	"github.com/jawher/mow.cli"
	"os"
	"sort"
	"time"
)

func start(path *string) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] MESSAGE"
		var (
			duration  = cmd.StringOpt("d duration", "25m", "duration of each stent")
			pomodoros = cmd.IntOpt("p pomodoros", 4, "number of pomodoros")
			message   = cmd.StringArg("MESSAGE", "", "descriptive name of the given task")
			tags      = cmd.StringsOpt("t tag", []string{}, "tags associated with this task")
		)
		cmd.Action = func() {
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := NewStore(*path)
			maybe(err)
			defer db.Close()
			task := &Task{
				Message:    *message,
				Tags:       *tags,
				NPomodoros: *pomodoros,
				Duration:   parsed,
			}
			runner, err := NewTaskRunner(task, db)
			maybe(err)
			maybe(runner.Run())
		}
	}
}

func initialize(path *string) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			db, err := NewStore(*path)
			maybe(err)
			defer db.Close()
			maybe(initDB(db))
		}
	}
}

func list(path *string) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			asJSON = cmd.BoolOpt("json", false, "output task history as JSON")
			assend = cmd.BoolOpt("assend", true, "sort tasks assending in age")
			limit  = cmd.IntOpt("n limit", 0, "limit the number of results by n")
		)
		cmd.Action = func() {
			db, err := NewStore(*path)
			maybe(err)
			defer db.Close()
			tasks, err := db.ReadTasks()
			maybe(err)
			if *assend {
				sort.Sort(sort.Reverse(ByID(tasks)))
			}
			if *limit > 0 && (len(tasks) > *limit) {
				tasks = tasks[0:*limit]
			}
			if *asJSON {
				maybe(json.NewEncoder(os.Stdout).Encode(tasks))
				return
			}
			config, _ := NewConfig(*path + "/config.json")
			summerizeTasks(config, tasks)
		}
	}
}

func _delete(path *string) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] TASK_ID"
		var taskID = cmd.IntArg("TASK_ID", -1, "task to delete")
		cmd.Action = func() {
			db, err := NewStore(*path)
			maybe(err)
			defer db.Close()
			maybe(db.DeleteTask(*taskID))
		}
	}
}

func main() {
	app := cli.App("pomo", "Pomodoro CLI")
	app.Spec = "[OPTIONS]"
	var (
		path = app.StringOpt("p path", defaultConfigPath(), "path to the pomo config directory")
	)
	app.Version("v version", Version)
	app.Command("start s", "start a new task", start(path))
	app.Command("init", "initialize the sqlite database", initialize(path))
	app.Command("list l", "list historical tasks", list(path))
	app.Command("delete d", "delete a stored task", _delete(path))
	app.Run(os.Args)
}
