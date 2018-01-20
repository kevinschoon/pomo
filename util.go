package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/user"
	"time"
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
		fmt.Printf("%d: [%s] ", task.ID, task.Duration.Truncate(time.Second))
		// a list of green/red pomodoros
		// green[x x] red[x x]
		fmt.Printf("[")
		for i := 0; i < task.NPomodoros; i++ {
			if i > 0 {
				fmt.Printf(" ")
			}
			if len(task.Pomodoros) >= i {
				color.New(color.FgGreen).Printf("X")
			} else {
				color.New(color.FgRed).Printf("X")
			}
		}
		fmt.Printf("]")
		if len(task.Tags) > 0 {
			fmt.Printf(" [")
			for i, tag := range task.Tags {
				if color, ok := config.Colors[tag]; ok {
					if i > 0 {
						fmt.Printf(" ")
					}
					color.Printf("%s", tag)
				}
			}
			fmt.Printf("]")
		}
		fmt.Printf(" - %s", task.Message)
		fmt.Printf("\n")
	}
}
