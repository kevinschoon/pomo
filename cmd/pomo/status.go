package main

import (
	"os"

	cli "github.com/jawher/mow.cli"

	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/runner"
	"github.com/kevinschoon/pomo/pkg/runner/client"
)

func status(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			c, err := client.NewSocketClient(cfg.SocketPath)
			if err != nil {
				status := runner.Status{}
				status.Write(os.Stdout)
				return
			}
			defer c.Close()
			status, err := c.Status()
			maybe(err)
			status.Write(os.Stdout)
		}
	}
}
