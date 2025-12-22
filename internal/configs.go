package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HomeDir     string
	CoursesDir  string
	CacheDir    string
	VideoExt    string
	DatabaseURL string
}

func LoadConfig() (*Config, error) {
	initEnv()
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
		HomeDir:     filepath.Clean(home),
		CoursesDir:  filepath.Join(home, "Courses"),
		CacheDir:    cache,
		VideoExt:    ".mkv",
		DatabaseURL: viper.GetString("database.url"),
	}

	return cfg, nil
}

func initEnv() {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	_ = viper.BindEnv("database.url", "DATABASE_URL")

	viper.SetDefault("database.url", "db.sqlite")
}
