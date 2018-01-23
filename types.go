package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"time"

	"github.com/kevinschoon/pomo/libnotify"
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
	Wheel           *Wheel
}

// Wheel keeps track of an ASCII spinner
type Wheel struct {
	state int
}

func (w *Wheel) String() string {
	switch w.state {
	case 0:
		w.state++
		return "|"
	case 1:
		w.state++
		return "/"
	case 2:
		w.state++
		return "-"
	case 3:
		w.state = 0
		return "\\"
	}
	return ""
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
	ID      int    `json:"id"`
	Message string `json:"message"`
	// Array of completed pomodoros
	Pomodoros []*Pomodoro `json:"pomodoros"`
	// Free-form tags associated with this task
	Tags []string `json:"tags"`
	// Number of pomodoros for this task
	NPomodoros int `json:"n_pomodoros"`
	// Duration of each pomodoro
	Duration time.Duration `json:"duration"`
}

// ByID is a sortable array of tasks
type ByID []*Task

func (b ByID) Len() int           { return len(b) }
func (b ByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByID) Less(i, j int) bool { return b[i].ID < b[j].ID }

// Pomodoro is a unit of time to spend working
// on a single task.
type Pomodoro struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Duration returns the runtime of the pomodoro
func (p Pomodoro) Duration() time.Duration {
	return (p.End.Sub(p.Start))
}

// Notifier implements a system specific
// notification. On Linux this libnotify.
// TODO: OSX, Windows(?)
type Notifier interface {
	Begin(int, Task) error
	Break(Task) error
	Finish(Task) error
}

// LibNotifier implements a Linux
// notifier client.
type LibNotifier struct {
	client   *libnotify.Client
	iconPath string
}

func NewLibNotifier() Notifier {
	ln := &LibNotifier{
		client: libnotify.NewClient(),
	}
	// Write the tomato icon to a temp path
	raw := MustAsset("tomato-icon.png")
	fp, _ := ioutil.TempFile("", "pomo")
	fp.Write(raw)
	ln.iconPath = fp.Name()
	fp.Close()
	return ln
}

func (ln LibNotifier) Begin(count int, t Task) error {
	return ln.client.Notify(libnotify.Notification{
		Title: t.Message,
		Body:  fmt.Sprintf("Task is starting (%d/%d pomodoros)", count, t.NPomodoros),
		Icon:  ln.iconPath,
	})
}

func (ln LibNotifier) Break(t Task) error {
	return ln.client.Notify(libnotify.Notification{
		Title:   t.Message,
		Urgency: "critical",
		Body:    fmt.Sprintf("Time to take a break!\nPress enter at the console to initiate the break."),
		Icon:    ln.iconPath,
	})
}

func (ln LibNotifier) Finish(t Task) error {
	return ln.client.Notify(libnotify.Notification{
		Title:   t.Message,
		Urgency: "critical",
		Body:    fmt.Sprintf("This task session is complete!"),
		Icon:    ln.iconPath,
	})
}
