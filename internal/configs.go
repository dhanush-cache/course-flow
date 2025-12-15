package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	HomeDir    string
	CoursesDir string
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		HomeDir:    filepath.Clean(home),
		CoursesDir: filepath.Join(home, "Courses"),
	}

	return cfg, nil
}
