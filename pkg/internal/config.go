package pomo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
)

const (
	defaultDateTimeFmt = "2006-01-02 15:04"
)

// Config represents user preferences
type Config struct {
	Colors      *ColorMap `json:"colors"`
	DateTimeFmt string    `json:"dateTimeFmt"`
	BasePath    string    `json:"basePath"`
	DBPath      string    `json:"dbPath"`
	SocketPath  string    `json:"socketPath"`
	IconPath    string    `json:"iconPath"`
	// Publish pushes updates to the configured
	// SocketPath rather than listening for requests
	Publish bool `json:"publish"`
	// PublishJson pushes socket updates as a JSON
	// encoded status message instead of string formatted
	PublishJson bool `json:"publishJson"`
	// If Publish is true, provide a socket path to publish to
	PublishSocketPath string `json:"publishSocketPath"`
}

type ColorMap struct {
	colors map[string]*color.Color
	tags   map[string]string
}

func (c *ColorMap) Get(name string) *color.Color {
	if color, ok := c.colors[name]; ok {
		return color
	}
	return nil
}

func (c *ColorMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.tags)
}

func (c *ColorMap) UnmarshalJSON(raw []byte) error {
	lookup := map[string]*color.Color{
		"black":     color.New(color.FgBlack),
		"hiblack":   color.New(color.FgHiBlack),
		"blue":      color.New(color.FgBlue),
		"hiblue":    color.New(color.FgHiBlue),
		"cyan":      color.New(color.FgCyan),
		"hicyan":    color.New(color.FgHiCyan),
		"green":     color.New(color.FgGreen),
		"higreen":   color.New(color.FgHiGreen),
		"magenta":   color.New(color.FgMagenta),
		"himagenta": color.New(color.FgHiMagenta),
		"red":       color.New(color.FgRed),
		"hired":     color.New(color.FgHiRed),
		"white":     color.New(color.FgWhite),
		"hiwrite":   color.New(color.FgHiWhite),
		"yellow":    color.New(color.FgYellow),
		"hiyellow":  color.New(color.FgHiYellow),
	}
	cm := &ColorMap{
		colors: map[string]*color.Color{},
		tags:   map[string]string{},
	}
	err := json.Unmarshal(raw, &cm.tags)
	if err != nil {
		return err
	}
	for tag, colorName := range cm.tags {
		if color, ok := lookup[colorName]; ok {
			cm.colors[tag] = color
		}
	}
	*c = *cm
	return nil
}

func LoadConfig(configPath string, config *Config) error {
	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		os.MkdirAll(path.Dir(configPath), 0755)
		// Create an empty config file
		// if it does not already exist.
		if os.IsNotExist(err) {
			raw, _ := json.Marshal(map[string]string{})
			err := ioutil.WriteFile(configPath, raw, 0644)
			if err != nil {
				return err
			}
			return LoadConfig(configPath, config)
		}
		return err
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		return err
	}
	if config.DateTimeFmt == "" {
		config.DateTimeFmt = defaultDateTimeFmt
	}
	if config.BasePath == "" {
		config.BasePath = path.Dir(configPath)
		err := os.MkdirAll(config.BasePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if config.DBPath == "" {
		config.DBPath = path.Join(xdg.DataHome, "pomo", "pomo.db")
		err := os.MkdirAll(path.Dir(config.DBPath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	if config.SocketPath == "" {
		config.SocketPath = path.Join(xdg.RuntimeDir, "pomo.sock")
		err := os.MkdirAll(path.Dir(config.SocketPath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	if config.IconPath == "" {
		config.IconPath = path.Join(xdg.DataHome, "pomo", "icon.png")
		err := os.MkdirAll(path.Dir(config.IconPath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	if config.Publish && (config.PublishSocketPath == "" || config.PublishSocketPath == config.SocketPath) {
		return fmt.Errorf("'publish' option now requires 'publishSocketPath' which must not be the same as 'socketPath'")
	}

	return nil
}
