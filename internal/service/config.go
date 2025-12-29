package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"text/tabwriter"

	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/db"
	"github.com/dhanush-cache/course-flow/internal/mosh"
)

type CourseWithPlatform struct {
	ID            string
	Slug          string
	Name          string
	PlatformTitle string
}

func GetCoursesWithPlatform(
	ctx context.Context,
	queries *db.Queries,
	cfg *config.Config,
) ([]CourseWithPlatform, error) {
	rows, err := queries.ListCoursesWithPlatforms(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]CourseWithPlatform, 0, len(rows))

	for _, r := range rows {
		course, err := mosh.CourseCache(mosh.GetData, cfg)(r.Slug, cfg)
		if err != nil {
			return nil, err
		}

		result = append(result, CourseWithPlatform{
			ID:            r.ID,
			Slug:          r.Slug,
			Name:          course.Name,
			PlatformTitle: r.PlatformTitle,
		})
	}

	return result, nil
}

func ListConfigs(cfg *config.Config) error {
	ctx := context.Background()

	dbDriver, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer dbDriver.Close()

	queries := db.New(dbDriver)

	courses, err := GetCoursesWithPlatform(ctx, queries, cfg)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, c := range courses {
		line := fmt.Sprintf(
			"%s\t%s\t%s\n",
			c.ID,
			c.Name,
			c.PlatformTitle,
		)
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
	}

	return w.Flush()
}
