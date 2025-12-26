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

func ListConfigs(cfg *config.Config) error {
	ctx := context.Background()
	dbDriver, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	queries := db.New(dbDriver)

	configs, err := queries.ListCoursesWithPlatforms(ctx)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, c := range configs {
		course, err := mosh.CourseCache(mosh.GetData, cfg)(c.Slug, cfg)
		if err != nil {
			return err
		}
		line := fmt.Sprintf(
			"%s\t%s\t%s\n",
			c.ID,
			course.Name,
			c.PlatformTitle,
		)
		_, err = w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}
