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
	TickTime           = 100 * time.Millisecond
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

func DefaultConfig() *Config {
	return &Config{
		DateTimeFmt: defaultDateTimeFmt,
		DBPath:      DefaultSharePath() + "/pomo.db",
		SocketPath:  DefaultSharePath() + "/pomo.sock",
		IconPath:    DefaultSharePath() + "/pomo.png",
	}
}

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

func DefaultSharePath() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(u.HomeDir, "/.local/share/pomo")
}

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
