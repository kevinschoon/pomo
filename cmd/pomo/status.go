package main

import (
	"fmt"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"

	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/runner/client"
	"github.com/kevinschoon/pomo/pkg/ui"
)

func status(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			watch    = cmd.BoolOpt("w watch", false, "continuously query for status")
			duration = cmd.StringOpt("d duration", "900ms", "watch interval")
			tmpl     = cmd.StringOpt("t tmpl", ui.DefaultStatusTmpl, "status go template")
		)
		cmd.Action = func() {
			duration, err := time.ParseDuration(*duration)
			maybe(err)
			if *watch {
				wheel := ui.Wheel(0)
				for {
					time.Sleep(duration)
					cli, err := client.NewSocketClient(cfg.SocketPath)
					if err != nil {
                        fmt.Println(ui.TemplateStatus(nil, nil, *tmpl))
						continue
					}
					status, err := cli.Status()
					if err != nil {
						continue
					}
                    fmt.Println(ui.TemplateStatus(status, &wheel, *tmpl))
				}
			} else {
				cli, err := client.NewSocketClient(cfg.SocketPath)
				if err != nil {
					fmt.Println(ui.TemplateStatus(nil, nil, *tmpl))
					os.Exit(1)
				}

				defer cli.Close()
				status, err := cli.Status()
				if err != nil {
					fmt.Println(ui.TemplateStatus(nil, nil, *tmpl))
					os.Exit(1)
				}
				fmt.Println(ui.TemplateStatus(status, nil, *tmpl))
			}
		}
	}
}
