package main

import (
	"fmt"
	"os/exec"
	"time"
)

// Task describes some activity
type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	count    int
	duration time.Duration
}

// Record is a stetch of work performed on a
// specific task.
type Record struct {
	TaskID int       `json:"task_id"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

// Prompter prompts a user with a message.
type Prompter interface {
	Prompt(string) error
}

// I3 implements a prompter for i3
type I3 struct{}

func (i *I3) Prompt(message string) error {
	raw, err := exec.Command(
		"/bin/i3-nagbar",
		"-m",
		message,
	).Output()
	if err != nil {
		return err
	}
	fmt.Println(string(raw))
	return nil
}
