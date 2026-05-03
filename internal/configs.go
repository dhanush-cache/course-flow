package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	CodeWithMosh = "codewithmosh"
	DreamsOfCode = "dreamsofcode"
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

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	cache := filepath.Join(cacheDir, "course-flow")

	if err := os.MkdirAll(cache, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create cache dir: %w", err)
	}

	dbURL := viper.GetString("database.url")
	if dbURL == "" {
		dbURL = filepath.Join(home, ".local", "share", "course-flow", "db.sqlite")
	}

	dbDir := filepath.Dir(dbURL)
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create database dir: %w", err)
	}

	cfg := &Config{
		HomeDir:     filepath.Clean(home),
		CoursesDir:  filepath.Join(home, "Courses"),
		CacheDir:    cache,
		VideoExt:    ".mkv",
		DatabaseURL: dbURL,
	}

	return cfg, nil
}

func initEnv() {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	_ = viper.BindEnv("database.url", "DATABASE_URL")
}
