package main

import (
	"os"

	cli "github.com/jawher/mow.cli"

	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/version"
)

func main() {
	app := cli.App("pomo", "Pomodoro CLI")
	app.LongDesc = "Pomo helps you track what you did, how long it took you to do it, and how much effort you expect it to take."
	app.Spec = "[OPTIONS]"
	var (
		cfg    = config.DefaultConfig()
		path   = app.StringOpt("p path", config.DefaultConfigPath(), "path to the pomo config directory")
		asJSON = app.BoolOpt("json", false, "output as json")
	)
	app.Before = func() {
		maybe(config.Load(*path, cfg))
		cfg.JSON = *asJSON
	}
	app.Version("v version", version.Version)
	app.Command("start s", "start a new task", start(cfg))
	app.Command("init", "initialize the sqlite database", initialize(cfg))
	app.Command("config cf", "display the current configuration", getConfig(cfg))
	app.Command("create c", "create a new task without starting", create(cfg))
	app.Command("edit e", "edit a resource", edit(cfg))
	app.Command("begin b", "begin requested pomodoro", begin(cfg))
	app.Command("get g", "get one or more tasks", get(cfg))
	app.Command("delete d", "delete a resource", deleteTask(cfg))
	app.Command("status st", "output the current status", status(cfg))
	app.Run(os.Args)
}
