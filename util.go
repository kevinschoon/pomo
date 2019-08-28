package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/user"
	"path"
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
