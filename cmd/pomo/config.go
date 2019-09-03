package main

import (
	"encoding/json"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/kevinschoon/pomo/pkg/config"
)

func getConfig(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(cfg))
		}
	}
}
