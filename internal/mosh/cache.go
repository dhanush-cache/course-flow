package mosh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	config "github.com/dhanush-cache/course-flow/internal"
)

type TokenFetchFunc func() (*Token, error)

type CourseFetchFunc func(slug string) (*Course, error)

type CoursesFetchFunc func() (*[]Course, error)

// TokenCache caches the result of a TokenFetchFunc based on the provided key.
func TokenCache(fn TokenFetchFunc, key string) TokenFetchFunc {
	return func() (*Token, error) {
		return cacheResult(
			key,
			fn,
			func(t *Token) bool {
				return t.ExpiresAt.After(time.Now())
			},
		)
	}
}

// CourseCache caches the result of a CourseFetchFunc based on the course slug.
func CourseCache(fn CourseFetchFunc) CourseFetchFunc {
	return func(slug string) (*Course, error) {
		return cacheResult(
			slug,
			func() (*Course, error) {
				return fn(slug)
			},
			nil,
		)
	}
}

// CoursesCache caches the result of a CoursesFetchFunc based on the provided key.
func CoursesCache(fn CoursesFetchFunc, key string) CoursesFetchFunc {
	return func() (*[]Course, error) {
		return cacheResult(
			key,
			func() (*[]Course, error) {
				return fn()
			},
			nil,
		)
	}
}

// cacheResult is a generic function that handles caching logic.
func cacheResult[T any](
	cacheKey string,
	fetch func() (*T, error),
	isValid func(*T) bool,
) (*T, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	cacheFile := filepath.Join(cfg.CacheDir, cacheKey+".json")

	if data, err := os.ReadFile(cacheFile); err == nil {
		var cached T
		if err := json.Unmarshal(data, &cached); err == nil {
			if isValid == nil || isValid(&cached) {
				// TODO: remove this debugging statement
				fmt.Println("CourseCache hit:", cacheFile)
				return &cached, nil
			}
		}
	}

	result, err := fetch()
	if err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write cache file: %w", err)
	}

	return result, nil
}
