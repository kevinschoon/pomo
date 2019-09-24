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
	defaultDateTimeFmt = "2006-01-02 15:04"
	// TickTime is the internal refresh time
	// used across the UI and internal timers
	TickTime = 100 * time.Millisecond
)

// Config represents user preferences
type Config struct {
	Colors         *ColorMap `json:"colors"`
	CurrentProject int64     `json:"currentProject"`
	JSON           bool      `json:"json"`
	DateTimeFmt    string    `json:"dateTimeFmt"`
	DBPath         string    `json:"dbPath"`
	SocketPath     string    `json:"socketPath"`
	IconPath       string    `json:"iconPath"`
}

// DefaultConfig returns the default Pomo configuration
func DefaultConfig() *Config {
	return &Config{
		DateTimeFmt: defaultDateTimeFmt,
		DBPath:      DefaultSharePath() + "/pomo.db",
		SocketPath:  DefaultSharePath() + "/pomo.sock",
		IconPath:    DefaultSharePath() + "/pomo.png",
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
