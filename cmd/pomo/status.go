package main

import (
	"fmt"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"

	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/runner"
	"github.com/kevinschoon/pomo/pkg/runner/client"
	"github.com/kevinschoon/pomo/pkg/ui"
)

const tomato rune = 0x1F345

func printStatus(path string, wheel *ui.Wheel) error {
	c, err := client.NewSocketClient(path)
	if err != nil {
		fmt.Printf("%c -\n", tomato)
		return err
	}
	defer c.Close()
	status, err := c.Status()
	if err != nil {
		fmt.Printf("%c -\n", tomato)
		return err
	}
	state := string(status.State.String()[0])
	if status.State == runner.RUNNING {
		remaining := (status.Duration - status.TimeRunning).Truncate(time.Second)
		fmt.Printf("%c %s [%d/%d] %s %s\n", tomato, state, status.Count, status.NPomodoros, wheel.String(), remaining)
	} else if status.State == runner.SUSPENDED {
		suspended := status.TimeSuspended.Truncate(time.Second)
		fmt.Printf("%c %s [%d/%d] %s +%s\n", tomato, state, status.Count, status.NPomodoros, wheel.Reverse(), suspended)
	} else {
		fmt.Printf("%c %s [%d/%d] -\n", tomato, state, status.Count, status.NPomodoros)
	}
	return nil
}

func status(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			watch    = cmd.BoolOpt("w watch", false, "continuously query for status")
			duration = cmd.StringOpt("d duration", "900ms", "watch interval")
		)
		cmd.Action = func() {
			duration, err := time.ParseDuration(*duration)
			maybe(err)
			if *watch {
				wheel := ui.Wheel(0)
				for {
					printStatus(cfg.SocketPath, &wheel)
					time.Sleep(duration)
				}
			} else {
				wheel := ui.Wheel(2)
				err := printStatus(cfg.SocketPath, &wheel)
				if err != nil {
					os.Exit(1)
				}
			}
		}
	}
}
