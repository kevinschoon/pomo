package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/fatih/color"
)

const (
	defaultDateTimeFmt = "2006-01-02 15:04"
)

// Config represents user preferences
type Config struct {
	Colors      *ColorMap `json:"colors"`
	DateTimeFmt string    `json:"dateTimeFmt"`
	DBPath      string    `json:"dbPath"`
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

func NewConfig(configPath string) (*Config, error) {
	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		// Create an empty config file
		// if it does not already exist.
		if os.IsNotExist(err) {
			raw, _ := json.Marshal(map[string]string{})
			err := ioutil.WriteFile(configPath, raw, 0644)
			if err != nil {
				return nil, err
			}
			return NewConfig(configPath)
		}
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}
	if config.DateTimeFmt == "" {
		config.DateTimeFmt = defaultDateTimeFmt
	}
	if config.DBPath == "" {
		config.DBPath = path.Dir(configPath) + "/pomo.db"
	}
	return config, nil
}
