package main

import (
	cli "github.com/jawher/mow.cli"
)

func status(config *Config) func(*cli.Cmd) {
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
