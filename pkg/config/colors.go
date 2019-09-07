package config

import (
	"encoding/json"

	"github.com/fatih/color"
)

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
