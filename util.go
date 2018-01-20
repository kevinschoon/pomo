package main

import (
	"fmt"
	//"github.com/fatih/color"
	"os"
	"os/user"
)

func maybe(err error) {
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		os.Exit(1)
	}
}

func defaultConfigPath() string {
	u, err := user.Current()
	maybe(err)
	return u.HomeDir + "/.pomo"
}

func summerizeTasks(config *Config, tasks []*Task) {
	for _, task := range tasks {
		var tags string
		if len(task.Tags) > 0 {
			for i, tag := range task.Tags {
				if color, ok := config.Colors[tag]; ok {
					if i > 0 {
						tags += " "
					}
					tags += color.SprintfFunc()("%s", tag)
				}
			}
		}
		fmt.Printf("%d [%s]: %s\n", task.ID, tags, task.Message)
	}
}
