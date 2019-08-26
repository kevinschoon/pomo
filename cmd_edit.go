package main

/*

import (
	"fmt"
	"time"

	cli "github.com/jawher/mow.cli"
)

func edit(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] TASKID"
		var (
            taskID = cmd.IntArg("TASKID", 0, "task id")
			parent    = cmd.IntOpt("parent", 0, "parent task")
			duration  = cmd.StringOpt("d duration", "", "duration of each pomodoro")
			pomodoros = cmd.IntOpt("p pomodoros", 0, "number of pomodoros")
			message   = cmd.StringArg("MESSAGE", "", "descriptive name of the given task")
			tags      = cmd.StringsOpt("t tag", []string{}, "tags associated with this task")
		)
		cmd.Action = func() {
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
            task, err := ReadOne(store, int64(*taskID))
            maybe(err)
            var modified bool

            // changing parent task
            if *parent > 0 {

            }

			taskID, err := CreateOne(store,
				&Task{
					ParentID:  int64(*parent),
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
*/
