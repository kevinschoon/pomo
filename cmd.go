package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

func Cmd() {
	app := cli.App("pomo", "Pomodoro CLI")
	app.LongDesc = "Pomo helps you track what you did, how long it took you to do it, and how much effort you expect it to take."
	app.Spec = "[OPTIONS]"
	var (
		config = &Config{}
		path   = app.StringOpt("p path", defaultConfigPath(), "path to the pomo config directory")
		asJSON = app.BoolOpt("json", false, "output as json")
	)
	app.Before = func() {
		maybe(LoadConfig(*path, config))
		config.JSON = *asJSON
	}
	app.Version("v version", Version)
	app.Command("start s", "start a new task", start(config))
	app.Command("init", "initialize the sqlite database", initialize(config))
	app.Command("config cf", "display the current configuration", getConfig(config))
	app.Command("create c", "create a new task without starting", create(config))
	app.Command("edit e", "edit a resource", edit(config))
	app.Command("begin b", "begin requested pomodoro", begin(config))
	app.Command("get g", "get one or more tasks", get(config))
	app.Command("delete d", "delete a resource", _delete(config))
	app.Command("status st", "output the current status", status(config))
	app.Run(os.Args)
}
