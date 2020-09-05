package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	cli "github.com/jawher/mow.cli"
)

func start(config *Config) func(*cli.Cmd) {
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
			db, err := NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			task := &Task{
				Message:    *message,
				Tags:       *tags,
				NPomodoros: *pomodoros,
				Duration:   parsed,
			}
			maybe(db.With(func(tx *sql.Tx) error {
				id, err := db.CreateTask(tx, *task)
				if err != nil {
					return err
				}
				task.ID = id
				return nil
			}))
			runner, err := NewTaskRunner(task, config)
			maybe(err)
			server, err := NewServer(runner, config)
			maybe(err)
			server.Start()
			defer server.Stop()
			runner.Start()
			startUI(runner)
		}
	}
}

func create(config *Config) func(*cli.Cmd) {
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
			db, err := NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			task := &Task{
				Message:    *message,
				Tags:       *tags,
				NPomodoros: *pomodoros,
				Duration:   parsed,
			}
			maybe(db.With(func(tx *sql.Tx) error {
				taskId, err := db.CreateTask(tx, *task)
				if err != nil {
					return err
				}
				fmt.Println(taskId)
				return nil
			}))
		}
	}
}

func begin(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] TASK_ID"
		var (
			taskId = cmd.IntArg("TASK_ID", -1, "ID of Pomodoro to begin")
		)

		cmd.Action = func() {
			db, err := NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			var task *Task
			maybe(db.With(func(tx *sql.Tx) error {
				read, err := db.ReadTask(tx, *taskId)
				if err != nil {
					return err
				}
				task = read
				err = db.DeletePomodoros(tx, *taskId)
				if err != nil {
					return err
				}
				task.Pomodoros = []*Pomodoro{}
				return nil
			}))
			runner, err := NewTaskRunner(task, config)
			maybe(err)
			server, err := NewServer(runner, config)
			maybe(err)
			server.Start()
			defer server.Stop()
			runner.Start()
			startUI(runner)
		}
	}
}

func initialize(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			db, err := NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(initDB(db))
		}
	}
}

func list(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			asJSON   = cmd.BoolOpt("json", false, "output task history as JSON")
			assend   = cmd.BoolOpt("assend", false, "sort tasks assending in age")
			all      = cmd.BoolOpt("a all", true, "output all tasks")
			limit    = cmd.IntOpt("n limit", 0, "limit the number of results by n")
			duration = cmd.StringOpt("d duration", "24h", "show tasks within this duration")
		)
		cmd.Action = func() {
			duration, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(tx *sql.Tx) error {
				tasks, err := db.ReadTasks(tx)
				maybe(err)
				if *assend {
					sort.Sort(sort.Reverse(ByID(tasks)))
				}
				if !*all {
					tasks = After(time.Now().Add(-duration), tasks)
				}
				if *limit > 0 && (len(tasks) > *limit) {
					tasks = tasks[0:*limit]
				}
				if *asJSON {
					maybe(json.NewEncoder(os.Stdout).Encode(tasks))
					return nil
				}
				maybe(err)
				summerizeTasks(config, tasks)
				return nil
			}))
		}
	}
}

func _delete(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "TASK_IDS..."
		ids := cmd.IntsArg("TASK_IDS", nil, "tasks to delete")
		cmd.Action = func() {
			db, err := NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(tx *sql.Tx) error {
				return db.DeleteTasks(tx, *ids)
			}))
		}
	}
}

func _status(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			client, err := NewClient(config.SocketPath)
			if err != nil {
				outputStatus(Status{})
				return
			}
			defer client.Close()
			status, err := client.Status()
			maybe(err)
			outputStatus(*status)
		}
	}
}

func _config(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(config))
		}
	}
}

func main() {
	app := cli.App("pomo", "Pomodoro CLI")
	app.LongDesc = "Pomo helps you track what you did, how long it took you to do it, and how much effort you expect it to take."
	app.Spec = "[OPTIONS]"
	var (
		config = &Config{}
		path   = app.StringOpt("p path", defaultConfigPath(), "path to the pomo config directory")
	)
	app.Before = func() {
		maybe(LoadConfig(*path, config))
	}
	app.Version("v version", Version)
	app.Command("start s", "start a new task", start(config))
	app.Command("init", "initialize the sqlite database", initialize(config))
	app.Command("config cf", "display the current configuration", _config(config))
	app.Command("create c", "create a new task without starting", create(config))
	app.Command("begin b", "begin requested pomodoro", begin(config))
	app.Command("list l", "list historical tasks", list(config))
	app.Command("delete d", "delete a list of stored tasks", _delete(config))
	app.Command("status st", "output the current status", _status(config))
	app.Run(os.Args)
}
