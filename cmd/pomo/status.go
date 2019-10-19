package main

import (
	"fmt"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"

	"github.com/kevinschoon/pomo/pkg/config"
	stl "github.com/kevinschoon/pomo/pkg/display/template/status"
	"github.com/kevinschoon/pomo/pkg/runner/client"
)

func status(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			watch    = cmd.BoolOpt("w watch", false, "continuously query for status")
			duration = cmd.StringOpt("d duration", "900ms", "watch interval")
			tmpl     = cmd.StringOpt("t tmpl", stl.DefaultStatusTmpl, "status go template")
		)
		cmd.Action = func() {
			duration, err := time.ParseDuration(*duration)
			maybe(err)
			templater := stl.NewStatusBarTemplater(*tmpl)
			if *watch {
				for {
					time.Sleep(duration)
					cli, err := client.NewSocketClient(cfg.SocketPath)
					if err != nil {
						fmt.Println(templater(nil))
						continue
					}
					status, err := cli.Status()
					if err != nil {
						continue
					}
					fmt.Println(templater(status))
				}
			} else {
				cli, err := client.NewSocketClient(cfg.SocketPath)
				if err != nil {
					fmt.Println(templater(nil))
					os.Exit(1)
				}

				defer cli.Close()
				status, err := cli.Status()
				if err != nil {
					fmt.Println(templater(nil))
					os.Exit(1)
				}
				fmt.Println(templater(status))
			}
		}
	}
}
