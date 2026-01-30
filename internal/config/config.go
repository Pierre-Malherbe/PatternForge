package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNotFound is returned when the settings file does not exist.
var ErrNotFound = errors.New("settings file not found")

// Config holds the application settings.
type Config struct {
	SaveDirectory string `json:"save_directory"`
}

// DefaultSaveDirectory returns the default path for saving results.
func DefaultSaveDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "results"
	}
	return filepath.Join(home, "Documents", "PatternForge", "results")
}

// ConfigDir returns the configuration directory path.
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".patternforge"
	}
	return filepath.Join(home, ".config", "patternforge")
}

// ConfigPath returns the full path to the settings file.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "settings.json")
}

// Exists returns whether the settings file exists.
func Exists() bool {
	_, err := os.Stat(ConfigPath())
	return err == nil
}

// Load reads and unmarshals the settings file.
// Returns ErrNotFound if the file does not exist.
func Load() (Config, error) {
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, ErrNotFound
		}
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// Save marshals and writes the config to the settings file.
func Save(cfg Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(ConfigPath(), data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
