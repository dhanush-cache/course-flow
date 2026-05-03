package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/db"
	"github.com/dhanush-cache/course-flow/internal/platforms/dreamsofcode"
	"github.com/dhanush-cache/course-flow/internal/platforms/mosh"
)

type Options struct {
	CourseID  string
	ZipFiles  []string
	Directory string
}

func AddCourse(options Options, cfg *config.Config) error {
	ctx := context.Background()
	dbDriver, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	queries := db.New(dbDriver)

	course, err := queries.GetCourse(ctx, options.CourseID)
	if err != nil {
		return err
	}

	// TODO: Implement the url logic
	var fileNames []string
	if course.PlatformID == config.CodeWithMosh {
		fileNames, err = processCodeWithMosh(course, cfg)
	} else if course.PlatformID == config.DreamsOfCode {
		fileNames, err = processDreamsOfCode(course, cfg)
	}
	if err != nil {
		return err
	}
	if options.ZipFiles != nil {
		err = ExtractAndProcess(course.ID, options.ZipFiles, fileNames, cfg)
	} else if options.Directory != "" {
		err = Process(course.ID, options.Directory, fileNames, cfg)
	}

	return err
}

func ListCourses() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(cfg.CoursesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println(entry.Name())
		}
	}

	return nil
}

func processCodeWithMosh(course db.Courses, cfg *config.Config) ([]string, error) {
	data, err := mosh.CourseCache(mosh.GetData, cfg)(course.Slug, cfg)
	if err != nil {
		return nil, err
	}
	fileNames, err := mosh.GetFileNames(data, cfg)
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}

func processDreamsOfCode(course db.Courses, cfg *config.Config) ([]string, error) {
	cacheFile := filepath.Join(cfg.CacheDir, course.Slug+".json")
	fmt.Println(cacheFile)
	data, err := dreamsofcode.LoadCourse(cacheFile)
	if err != nil {
		return nil, err
	}
	fileNames, err := dreamsofcode.GetFileNames(data, cfg)
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}
