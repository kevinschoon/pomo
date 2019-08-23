package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/fatih/color"
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
	return path.Join(u.HomeDir, "/.pomo/config.json")
}

func makeUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func makeTimers(task Task) []*Timer {
	timers := make([]*Timer, len(task.Pomodoros))
	for i := 0; i < len(task.Pomodoros); i++ {
		runtime := task.Pomodoros[i].RunTime
		pauseTime := task.Pomodoros[i].PauseTime
		timers[i] = NewTimer(task.Duration, runtime, pauseTime)
	}
	return timers
}

// TODO: Make more configurable / possibly remove all together
func summerizeTasks(config *Config, tasks []*Task) {
	for _, task := range tasks {
		var start string
		if len(task.Pomodoros) > 0 {
			start = task.Pomodoros[0].Start.Format(config.DateTimeFmt)
		}
		fmt.Printf("%d: [%s] [%s] ", task.ID, start, task.Duration.Truncate(time.Second))
		// a list of green/yellow/red pomodoros
		// green indicates the pomodoro was finished normally
		// yellow indicates the break was exceeded by +5minutes
		// red indicates the pomodoro was never completed
		fmt.Printf("[")
		for i, pomodoro := range task.Pomodoros {
			if i > 0 {
				fmt.Printf(" ")
			}
			// pomdoro was paused more than 5m
			if pomodoro.PauseTime > task.Duration+5*time.Minute {
				color.New(color.FgYellow).Printf("X")
			} else {
				// pomodoro completed normally
				color.New(color.FgGreen).Printf("X")
			}
		}
		// each missed pomodoro
		/*
			for i := 0; i < task.NPomodoros-len(task.Pomodoros); i++ {
				if i > 0 || i == 0 && len(task.Pomodoros) > 0 {
					fmt.Printf(" ")
				}
				color.New(color.FgRed).Printf("X")
			}
		*/
		fmt.Printf("]")
		// Tags
		if len(task.Tags) > 0 {
			fmt.Printf(" [")
			for i, tag := range task.Tags {
				if i > 0 && i != len(task.Tags) {
					fmt.Printf(" ")
				}
				// user specified color mapping exists
				if color := config.Colors.Get(tag); color != nil {
					color.Printf("%s", tag)
				} else {
					// no color mapping
					fmt.Printf("%s", tag)
				}
			}
			fmt.Printf("]")
		}
		fmt.Printf(" - %s", task.Message)
		fmt.Printf("\n")
	}
}

func outputStatus(status Status) {
	state := "?"
	if status.State >= RUNNING {
		state = string(status.State.String()[0])
	}
	if status.State == RUNNING {
		fmt.Printf("%s [%d/%d]", state, status.Count, status.NPomodoros)
	} else {
		fmt.Printf("%s [%d/%d] -", state, status.Count, status.NPomodoros)
	}
}
