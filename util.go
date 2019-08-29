package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
)

func maybe(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func defaultConfigPath() string {
	u, err := user.Current()
	maybe(err)
	return path.Join(u.HomeDir, "/.pomo/config.json")
}

func parseTags(kvs []string) (map[string]string, error) {
	tags := map[string]string{}
	for _, kv := range kvs {
		split := strings.Split(kv, "=")
		if len(split) == 2 {
			tags[split[0]] = split[1]
		} else {
			return nil, fmt.Errorf("bad tag: %s", kv)
		}
	}
	return tags, nil
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

func promptConfirm(question string) error {
	reader := bufio.NewReader(os.Stdin)
	result, _ := reader.ReadString('\n')
	result = strings.Replace(result, "\n", "", -1)
	if result != question {
		return fmt.Errorf("cancelled")
	}
	return nil
}

func truncDuration(s string) string {
	if len(s) > 4 {
		return strings.Replace(strings.Replace(s, "0s", "", -1), "0m", "", -1)
	}
	return s
}
