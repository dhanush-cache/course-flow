package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	HomeDir    string
	CoursesDir string
	CacheDir   string
	VideoExt   string
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cache, err := filepath.Abs("data/cache")
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(cache, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create cache dir: %w", err)
	}

	cfg := &Config{
		HomeDir:    filepath.Clean(home),
		CoursesDir: filepath.Join(home, "Courses"),
		CacheDir:   cache,
		VideoExt:   ".mkv",
	}

	return cfg, nil
}
