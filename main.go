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
			db, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer db.Close()
			taskID, err := CreateOne(db, &Task{
				Message:   *message,
				Tags:      *tags,
				Pomodoros: NewPomodoros(*pomodoros),
				Duration:  parsed,
			})
			maybe(err)
			_, err = ReadOne(db, taskID)
			maybe(err)
			runner := NewTaskRunner()
			server, err := NewSocketServer(runner, config)
			maybe(err)
			// TODO catch error
			// runner.Start(task)
			go server.Serve()
			// defer server.Stop()
			client, err := NewSocketClient(config.SocketPath)
			maybe(err)
			maybe(startUI(client))
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
			db, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer db.Close()
			taskID, err := CreateOne(db,
				&Task{
					Message:   *message,
					Tags:      *tags,
					Pomodoros: NewPomodoros(*pomodoros),
					Duration:  parsed,
				})
			maybe(err)
			fmt.Printf("%d", taskID)
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
			db, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer db.Close()
			_, err = ReadOne(db, int64(*taskId))
			maybe(err)
			/*
				maybe(db.With(func(tx *sql.Tx) error {
					err := db.ReadTask(tx, task)
					if err != nil {
						return err
					}
					// TODO
					err = db.DeletePomodoros(tx, int64(*taskId), int64(-1))
					if err != nil {
						return err
					}
					task.Pomodoros = []*Pomodoro{}
					return nil
				}))
			*/
			runner := NewTaskRunner()
			maybe(err)
			server, err := NewSocketServer(runner, config)
			maybe(err)
			go server.Serve()
			client, err := NewSocketClient(config.SocketPath)
			maybe(err)
			maybe(startUI(client))
		}
	}
}

func initialize(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			db, err := NewSQLiteStore(config.DBPath)
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
			db, err := NewSQLiteStore(config.DBPath)
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
		cmd.Spec = "[OPTIONS] TASK_ID"
		var taskID = cmd.IntArg("TASK_ID", -1, "task to delete")
		cmd.Action = func() {
			db, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(tx *sql.Tx) error {
				return db.DeleteTask(tx, int64(*taskID))
			}))
		}
	}
}

func _status(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			client, err := NewSocketClient(config.SocketPath)
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
	app.Command("delete d", "delete a stored task", _delete(config))
	app.Command("status st", "output the current status", _status(config))
	app.Run(os.Args)
}
