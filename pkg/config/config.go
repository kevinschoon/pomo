package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"
)

const (
	// TickTime is the internal refresh time
	// used across the UI and internal timers
	TickTime = 100 * time.Millisecond
)

// Config represents user preferences
type Config struct {
	Colors         *ColorMap
	CurrentProject int64
	JSON           bool
	DBPath         string
	SocketPath     string
	IconPath       string
	// the number of snapshots to retain
	// if set to -1 snapshotting is disabled
	// if set to 0 all snapshots are retained
	Snapshots int
	// sets the default duration of pomodoros
	// when creating a new task
	DefaultDuration time.Duration
	// the default number of pomodoros that are configured
	// when creating a new task
	DefaultPomodoros int
}

type configJson struct {
	Colors         *ColorMap `json:"colors"`
	CurrentProject int64     `json:"current_project"`
	JSON           bool      `json:"json"`
	DBPath         string    `json:"db_path"`
	SocketPath     string    `json:"socket_path"`
	IconPath       string    `json:"icon_path"`
	// the number of snapshots to retain
	// if set to -1 snapshotting is disabled
	// if set to 0 all snapshots are retained
	Snapshots int `json:"history"`
	// sets the default duration of pomodoros
	// when creating a new task
	DefaultDuration string `json:"default_duration"`
	// the default number of pomodoros that are configured
	// when creating a new task
	DefaultPomodoros int `json:"default_pomodoros"`
}

func (c Config) MarshalJSON() ([]byte, error) {
	intermediate := configJson{
		Colors:           c.Colors,
		CurrentProject:   c.CurrentProject,
		JSON:             c.JSON,
		DBPath:           c.DBPath,
		SocketPath:       c.SocketPath,
		IconPath:         c.IconPath,
		Snapshots:        c.Snapshots,
		DefaultDuration:  c.DefaultDuration.String(),
		DefaultPomodoros: c.DefaultPomodoros,
	}
	return json.Marshal(intermediate)
}

func (c *Config) UnmarshalJSON(raw []byte) error {
	intermediate := &configJson{}
	err := json.Unmarshal(raw, intermediate)
	if err != nil {
		return err
	}
	c.Colors = intermediate.Colors
	c.CurrentProject = intermediate.CurrentProject
	c.JSON = intermediate.JSON
	c.SocketPath = intermediate.SocketPath
	c.IconPath = intermediate.IconPath
	c.Snapshots = intermediate.Snapshots
	if intermediate.DefaultDuration != "" {
		d, err := time.ParseDuration(intermediate.DefaultDuration)
		if err != nil {
			return err
		}
		c.DefaultDuration = d
	} else {
		c.DefaultDuration = 50 * time.Minute
	}
	c.DefaultPomodoros = intermediate.DefaultPomodoros
	return nil
}

// DefaultConfig returns the default Pomo configuration
func DefaultConfig() *Config {
	sharePath := DefaultSharePath()
	return &Config{
		DBPath:           path.Join(sharePath, "/pomo.db"),
		SocketPath:       path.Join(sharePath, "/pomo.sock"),
		IconPath:         path.Join(sharePath, "/pomo.png"),
		Snapshots:        10,
		DefaultDuration:  50 * time.Minute,
		DefaultPomodoros: 3,
	}
}

// GetConfigPath resolves the configuration path
// and checks if an alternate path has been
// specified via environment variable
func GetConfigPath() string {
	if os.Getenv("POMO_CONFIG_PATH") != "" {
		return os.Getenv("POMO_CONFIG_PATH")
	}
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(u.HomeDir, "/.config/pomo/config.json")
}

// DefaultSharePath returns the default path pomo
// stores it's SQLite and other data
func DefaultSharePath() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(u.HomeDir, "/.local/share/pomo")
}

// EnsurePaths ensures all needed paths have been
// created when initializing Pomo
func EnsurePaths(cfg *Config) error {
	_, err := os.Stat(path.Dir(cfg.DBPath))
	if os.IsNotExist(err) {
		return os.MkdirAll(path.Dir(cfg.DBPath), 0755)
	}
	_, err = os.Stat(path.Dir(cfg.IconPath))
	if os.IsNotExist(err) {
		return os.MkdirAll(path.Dir(cfg.IconPath), 0755)
	}
	_, err = os.Stat(path.Dir(cfg.SocketPath))
	if os.IsNotExist(err) {
		return os.MkdirAll(path.Dir(cfg.SocketPath), 0755)
	}
	return nil
}

// Load loads the given configuration from the path
func Load(configPath string, config *Config) error {
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
			return Load(configPath, config)
		}
		return err
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		return err
	}
	return nil
}
