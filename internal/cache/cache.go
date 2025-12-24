package cache

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Manager struct {
    Dir  string
    File string
}


// Creates a cache manager for the app
func New(app string, filename string) (*Manager, error) {
	// Check if home var exists
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		// Setting default folder as current
		dir := "cache"
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}

		return &Manager{
			Dir: dir,
			File: filepath.Join(dir, filename),
		}, nil

	}

	dir := filepath.Join(home, ".cache", app)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	return &Manager{
		Dir:  dir,
		File: filepath.Join(dir, filename),
	}, nil
}

// Checks if cache file exists
func (c *Manager) Exists() bool {
	_, err := os.Stat(c.File)
	return err == nil
}

// Checks if cache file exists and if its newer than maxAge
func (c *Manager) IsFresh(maxAge time.Duration) bool {
	info, err := os.Stat(c.File)
	if err != nil {
		return false
	}

	return time.Since(info.ModTime()) < maxAge
}

// Loads cache JSON into struct pointer
func (c *Manager) ReadJSON(v any) error {
	if !c.Exists() {
		return errors.New("cache does not exist")
	}

	data, err := os.ReadFile(c.File)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// Writes struct to cache file
func (c *Manager) WriteJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.File, data, 0o644)
}
