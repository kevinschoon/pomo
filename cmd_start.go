package main

import (
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
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			kvs, err := parseTags(*tags)
			maybe(err)
			task := &Task{
				Message:   *message,
				Tags:      kvs,
				Pomodoros: NewPomodoros(*pomodoros),
				Duration:  parsed,
			}
			maybe(store.With(func(s Store) error {
				err = store.CreateTask(task)
				if err != nil {
					return err
				}
				return store.ReadTask(task)
			}))
			server, err := NewSocketServer(task, store, config)
			maybe(err)
			shutdown := make(chan error)
			go func() {
				shutdown <- server.Serve()
			}()
			// runner.Start(task)
			// defer server.Stop()
			maybe(startUI(server))
			maybe(<-shutdown)
		}
	}

}
