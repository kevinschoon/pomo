package color

import (
	"encoding/json"

	"github.com/bcicen/color"
)

// TagPair describes the color configuration
// of a particular tag
type TagPair struct {
	Value string
	Color color.Color
	hex   string
}

type tagPairJSON struct {
	Value string `json:"value"`
	Color string `json:"color"`
	hex   string
}

func (t TagPair) MarshalJSON() ([]byte, error) {
	tp := &tagPairJSON{
		Value: t.Value,
		Color: t.hex,
	}
	return json.Marshal(tp)
}

func (t *TagPair) UnmarshalJSON(raw []byte) error {
	tp := &tagPairJSON{}
	err := json.Unmarshal(raw, tp)
	if err != nil {
		return err
	}
	t.Value = tp.Value
	c, err := color.NewHex(tp.Color)
	if err != nil {
		return err
	}
	t.Color = *c
	t.hex = tp.Color
	return nil
}

// Colors defines configurable colors for the Pomodoro CLI
type Colors struct {
	// color of progress below 50%
	Primary color.Color
	primary string
	//  color of progress above 50%
	Secondary color.Color
	secondary string
	// color of progress at 100%
	Tertiary color.Color
	tertiary string
	// map of tag key/value pair to colors
	Tags map[string]TagPair
}

type colorsJSON struct {
	Primary   string             `json:"primary"`
	Secondary string             `json:"secondary"`
	Tertiary  string             `json:"tertiary"`
	Tags      map[string]TagPair `json:"tags"`
}

// Get returns the color for a tag key/value pair
func (c Colors) Get(key, value string) color.Color {
	if cfg, ok := c.Tags[key]; ok {
		if cfg.Value == value {
			return cfg.Color
		}
	}
	return color.Color{}
}

// MarshalJSON marshals underlying tags
func (c *Colors) MarshalJSON() ([]byte, error) {
	cfg := &colorsJSON{
		Primary:   c.primary,
		Secondary: c.secondary,
		Tertiary:  c.tertiary,
		Tags:      c.Tags,
	}
	return json.Marshal(cfg)
}

// UnmarshalJSON returns a resolved ColorMap as JSON
func (c *Colors) UnmarshalJSON(raw []byte) error {
	cfg := &colorsJSON{}
	err := json.Unmarshal(raw, cfg)
	if err != nil {
		return err
	}
	c.Tags = cfg.Tags
	primary, err := color.NewHex(cfg.Primary)
	if err != nil {
		return err
	}
	c.Primary = *primary
	c.primary = cfg.Primary
	secondary, err := color.NewHex(cfg.Secondary)
	if err != nil {
		return err
	}
	c.Secondary = *secondary
	c.secondary = cfg.Secondary
	tertiary, err := color.NewHex(cfg.Tertiary)
	if err != nil {
		return err
	}
	c.Tertiary = *tertiary
	c.tertiary = cfg.Tertiary
	return nil
}
