package main

import (
	"encoding/json"
	"fmt"
	"github.com/jawher/mow.cli"
	"os"
	"time"
)

func maybe(err error) {
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		os.Exit(1)
	}
}

func start(path *string) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] MESSAGE"
		var (
			duration = cmd.StringOpt("d duration", "25m", "duration of each stent")
			stents   = cmd.IntOpt("s stents", 4, "number of working stents")
			message  = cmd.StringArg("MESSAGE", "", "descriptive name of the given task")
			tags     = cmd.StringsOpt("t tag", []string{}, "tags associated with this task")
		)
		cmd.Action = func() {
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := NewStore(*path)
			maybe(err)
			defer db.Close()
			task := Task{
				Message:  *message,
				Tags:     *tags,
				stents:   *stents,
				duration: parsed,
			}
			run(task, &I3{}, db)
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
		cmd.Action = func() {
			db, err := NewStore(*path)
			maybe(err)
			defer db.Close()
			tasks, err := db.ReadTasks()
			maybe(err)
			maybe(json.NewEncoder(os.Stdout).Encode(tasks))
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
		path = app.StringOpt("p path", defaultDBPath(), "path to the pomo state directory")
	)
	app.Command("start s", "start a new task", start(path))
	app.Command("init", "initialize the sqlite database", initialize(path))
	app.Command("list l", "list historical tasks", list(path))
	app.Command("delete d", "delete a stored task", _delete(path))
	app.Run(os.Args)
}
