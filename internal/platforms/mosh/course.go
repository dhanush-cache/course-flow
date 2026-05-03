package mosh

import (
	"fmt"
	"strings"

	"github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/utils"
)

// TODO: Improve the config management in this module

// GetData fetches course data for the given slug.
func GetData(slug string, cfg *config.Config) (*Course, error) {
	url := fmt.Sprintf("https://codewithmosh.com/p/%s", slug)

	var response CourseResponse
	err := ScrapePage(url, &response)
	if err != nil {
		return nil, err
	}

	course := &response.Props.PageProps.Course
	if course.Type == "bundle" {
		idMap := make(map[int]struct{})
		for _, id := range course.BundleContents {
			idMap[id] = struct{}{}
		}
		courses, err := CoursesCache(GetCourses, "courses", cfg)(cfg)
		if err != nil {
			return nil, err
		}
		for _, c := range *courses {
			if _, exists := idMap[c.ID]; exists {
				childCourse, err := CourseCache(GetData, cfg)(c.Slug, cfg)
				if err != nil {
					return nil, err
				}
				course.Courses = append(course.Courses, childCourse)
			}
		}
	}
	return course, nil
}

// GetCourses fetches the list of all courses.
func GetCourses(cfg *config.Config) (*[]Course, error) {
	url := "https://codewithmosh.com/courses"

	var response CoursesResponse
	err := ScrapePage(url, &response)
	if err != nil {
		return nil, err
	}

	return &response.Props.PageProps.Courses, nil
}

type ParentInfo struct {
	ParentName string
	Index      int
}

// GetFileNames retrieves file names for the given course.
func GetFileNames(course *Course, cfg *config.Config) ([]string, error) {
	return doGetFileNames(course, nil, cfg)
}

// doGetFileNames recursively retrieves file names for the given course.
func doGetFileNames(course *Course, parent *ParentInfo, cfg *config.Config) ([]string, error) {
	if course.Type == "bundle" {
		var result []string
		for i, childCourse := range course.Courses {
			filenames, err := doGetFileNames(childCourse, &ParentInfo{course.Name, i + 1}, cfg)
			if err != nil {
				return nil, err
			}
			result = append(result, filenames...)
		}
		return result, nil
	}
	var fileNames []string
	sections := course.Curriculum
	var name string
	if parent == nil {
		name = removeBuzzWords(course.Name)
	} else {
		name = fmt.Sprintf(
			"%s/Part - %d",
			removeBuzzWords(parent.ParentName),
			parent.Index,
		)
	}
	for s, section := range sections {
		for l, lesson := range section.Lessons {
			if lesson.IsVideo != 1 {
				continue
			}
			path := fmt.Sprintf(
				"%s/%02d - %s/%02d - %s%s",
				name,
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

// removeBuzzWords cleans up the course name by removing common buzz words.
func removeBuzzWords(name string) string {
	buzzWords := []string{
		"Mastering",
		"Mastery",
		"The Ultimate",
		"Ultimate",
		"The Complete",
		"Complete",
		"Series",
		"Bundle",
		"Crash Course",
		"Course",
		"for Beginners",
		".js",
	}
	for _, buzz := range buzzWords {
		name = strings.Replace(name, buzz, "", -1)
	}
	return strings.TrimSpace(name)
}
