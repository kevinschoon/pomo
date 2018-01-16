package main

import (
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

func startTask(task Task, prompter Prompter, db *Store) {
	taskID, err := db.AddTask(task)
	maybe(err)
	for i := 0; i < task.count; i++ {
		// Create a record for
		// this particular stent of work
		record := &Record{}
		// Prompt the client
		maybe(prompter.Prompt("Begin Working!"))
		record.Start = time.Now()
		// Wait the specified interval
		time.Sleep(task.duration)
		maybe(prompter.Prompt("Take a Break!"))
		// Record how long the user waited
		// until closing the notification
		record.End = time.Now()
		maybe(db.AddRecord(taskID, *record))
	}

}

func start(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS] NAME"
	var (
		duration = cmd.StringOpt("d duration", "25m", "duration of each stent")
		count    = cmd.IntOpt("c count", 4, "number of working stents")
		name     = cmd.StringArg("NAME", "", "descriptive name of the given task")
		path     = cmd.StringOpt("p path", defaultDBPath(), "path to the pomo state directory")
	)
	cmd.Action = func() {
		parsed, err := time.ParseDuration(*duration)
		maybe(err)
		db, err := NewStore(*path)
		maybe(err)
		defer db.Close()
		task := Task{
			Name:     *name,
			count:    *count,
			duration: parsed,
		}
		startTask(task, &I3{}, db)
	}
}

func initialize(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"
	var (
		path = cmd.StringOpt("p path", defaultDBPath(), "path to the pomo state directory")
	)
	cmd.Action = func() {
		db, err := NewStore(*path)
		maybe(err)
		defer db.Close()
		maybe(initDB(db))
	}
}

func list(cmd *cli.Cmd) {}

func main() {
	app := cli.App("pomo", "Pomodoro CLI")
	app.Spec = "[OPTIONS]"
	app.Command("start", "start a new task", start)
	app.Command("init", "initialize the sqlite database", initialize)
	app.Command("ls", "list historical tasks", list)
	app.Run(os.Args)
}
