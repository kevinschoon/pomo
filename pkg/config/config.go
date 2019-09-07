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
	Colors      *ColorMap `json:"colors"`
	JSON        bool      `json:"json"`
	DateTimeFmt string    `json:"dateTimeFmt"`
	BasePath    string    `json:"basePath"`
	DBPath      string    `json:"dbPath"`
	SocketPath  string    `json:"socketPath"`
	IconPath    string    `json:"iconPath"`
}

func DefaultConfig() *Config {
	return &Config{
		DateTimeFmt: defaultDateTimeFmt,
	}
}

func DefaultConfigPath() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(u.HomeDir, "/.pomo/config.json")
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
	if config.BasePath == "" {
		config.BasePath = path.Dir(configPath)
	}
	if config.DBPath == "" {
		config.DBPath = path.Join(config.BasePath, "/pomo.db")
	}
	if config.SocketPath == "" {
		config.SocketPath = path.Join(config.BasePath, "/pomo.sock")
	}
	if config.IconPath == "" {
		config.IconPath = path.Join(config.BasePath, "/icon.png")
	}
	return nil
}
