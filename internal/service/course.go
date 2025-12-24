package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/db"
	"github.com/dhanush-cache/course-flow/internal/mosh"
)

func AddCourse(key string, zipFiles []string, cfg *config.Config) error {
	ctx := context.Background()
	dbDriver, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	queries := db.New(dbDriver)

	course, err := queries.GetCourse(ctx, key)
	if err != nil {
		return err
	}

	// TODO: Implement the url logic
	if course.PlatformID == config.CodeWithMosh {
		err = processCodeWithMosh(course, zipFiles, cfg)
	}
	if err != nil {
		return err
	}
	return nil
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

func processCodeWithMosh(course db.Courses, zipFiles []string, cfg *config.Config) error {
	data, err := mosh.CourseCache(mosh.GetData, cfg)(course.Slug, cfg)
	if err != nil {
		return err
	}
	fileNames, err := mosh.GetFileNames(data, cfg)
	if err != nil {
		return err
	}
	err = Process(zipFiles, fileNames, cfg)
	if err != nil {
		return err
	}
	return nil
}
