package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"

	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/runner"
	"github.com/kevinschoon/pomo/pkg/runner/client"
)

const tomato rune = 0x1F345

func status(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			c, err := client.NewSocketClient(cfg.SocketPath)
			if err != nil {
				status := runner.Status{}
				fmt.Print(status.String())
				return
			}
			defer c.Close()
			status, err := c.Status()
			maybe(err)
			fmt.Printf("%c %s", tomato, status.String())
		}
	}
}
