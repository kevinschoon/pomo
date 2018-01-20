package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os/exec"
	"time"
)

// RefreshInterval is the frequency at which
// the display is updated.
const RefreshInterval = 800 * time.Millisecond

// Message is used internally for updating
// the display.
type Message struct {
	Start           time.Time
	Duration        time.Duration
	Pomodoros       int
	CurrentPomodoro int
}

// Config represents user preferences
type Config struct {
	Colors map[string]*color.Color
}

var colorMap = map[string]*color.Color{
	"red":   color.New(color.FgRed),
	"blue":  color.New(color.FgBlue),
	"green": color.New(color.FgGreen),
	"white": color.New(color.FgWhite),
}

func (c *Config) UnmarshalJSON(raw []byte) error {
	config := &struct {
		Colors map[string]string `json:"colors"`
	}{}
	err := json.Unmarshal(raw, config)
	if err != nil {
		return err
	}
	for key, name := range config.Colors {
		if color, ok := colorMap[name]; ok {
			c.Colors[key] = color
		} else {
			return fmt.Errorf("bad color choice: %s", name)
		}
	}
	return nil
}

func NewConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{
		Colors: map[string]*color.Color{},
	}
	return config, json.Unmarshal(raw, config)
}

// Task describes some activity
type Task struct {
	ID        int         `json:"id"`
	Message   string      `json:"message"`
	Pomodoros []*Pomodoro `json:"pomodoros"`
	// Free-form tags associated with this task
	Tags []string `json:"tags"`
	// Number of pomodoros for this task
	pomodoros int
	duration  time.Duration
}

// Pomodoro is a unit of time to spend working
// on a single task.
type Pomodoro struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Prompter prompts a user with a message.
type Prompter interface {
	Prompt(string) error
}

// I3 implements a prompter for i3
type I3 struct{}

func (i *I3) Prompt(message string) error {
	_, err := exec.Command(
		"/bin/i3-nagbar",
		"-m",
		message,
	).Output()
	if err != nil {
		return err
	}
	return nil
}
