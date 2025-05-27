package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// Config represents the optional configuration file structure.
type Config struct {
	MAC string `json:"mac"`
}

// Load reads ~/.config/embermug.json if it exists.
func Load() (*Config, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "embermug.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// DefaultMAC returns the MAC address from the config file, if present.
func DefaultMAC() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}
	return cfg.MAC, nil
}
