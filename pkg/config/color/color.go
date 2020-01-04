package color

import (
	"encoding/json"

	"github.com/bcicen/color"
)

func DefaultColors() *Colors {
	return &Colors{
		ConfigSrc: ConfigSrc{
			PrimaryStr:   "#FF6985",
			SecondaryStr: "#FFEC42",
			TertiaryStr:  "#7FE9A2",
			Tags:         make(map[string]TagPair),
		},
		Primary:   color.MustNewHex("#FF6985"),
		Secondary: color.MustNewHex("#FFEC42"),
		Tertiary:  color.MustNewHex("#7FE9A2"),
	}
}

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
	ConfigSrc

	Primary   *color.Color
	Secondary *color.Color
	Tertiary  *color.Color
}

type ConfigSrc struct {
	PrimaryStr   string             `json:"primary"`   // color of progress below 50%
	SecondaryStr string             `json:"secondary"` // color of progress above 50%
	TertiaryStr  string             `json:"tertiary"`  // color of progress at 100%
	Tags         map[string]TagPair `json:"tags"`      // map of tag key/value pair to colors
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
	return json.Marshal(&c.ConfigSrc)
}

// UnmarshalJSON returns a resolved ColorMap as JSON
func (c *Colors) UnmarshalJSON(raw []byte) error {
	var cfg ConfigSrc
	var err error

	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		return err
	}

	c.ConfigSrc = cfg

	c.Primary, err = color.NewHex(c.PrimaryStr)
	if err != nil {
		return err
	}
	c.Secondary, err = color.NewHex(c.SecondaryStr)
	if err != nil {
		return err
	}
	c.Tertiary, err = color.NewHex(c.TertiaryStr)
	if err != nil {
		return err
	}

	return nil
}
