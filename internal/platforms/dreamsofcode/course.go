package dreamsofcode

import (
	"fmt"

	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/utils"
)

// GetFileNames retrieves file names for the given course.
func GetFileNames(course *Course, cfg *config.Config) ([]string, error) {
	var fileNames []string

	for s, section := range course.Sections {
		for l, lesson := range section.Lessons {
			path := fmt.Sprintf(
				"%s/%02d - %s/%02d - %s%s",
				course.Name,
				s+1,
				section.Name,
				l+1,
				lesson.Name,
				cfg.VideoExt,
			)

			fileNames = append(fileNames, utils.CleanPath(path))
		}
	}

	return fileNames, nil
}
