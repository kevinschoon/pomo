package main

import (
	cli "github.com/jawher/mow.cli"
)

func begin(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] TASK_ID"
		var (
			taskId = cmd.IntArg("TASK_ID", -1, "ID of Pomodoro to begin")
		)

		cmd.Action = func() {
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			task, err := ReadOne(store, int64(*taskId))
			maybe(err)
			server, err := NewSocketServer(task, store, config)
			maybe(err)
			go server.Serve()
			client, err := NewSocketClient(config.SocketPath)
			maybe(err)
			maybe(startUI(client))
		}
	}
}
