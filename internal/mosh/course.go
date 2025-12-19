package mosh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/utils"
)

// GetData fetches course data for the given slug.
func GetData(slug string) (*Course, error) {
	token, err := TokenCache(GetToken, "token")()
	if err != nil {
		return nil, fmt.Errorf("error getting token: %v", err)
	}

	url := fmt.Sprintf("https://codewithmosh.com/_next/data/%s/p/%s.json", token.Value, slug)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Printf("Making GET request to: %s\n", url)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response CourseResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	course := &response.PageProps.Course
	if course.Type == "bundle" {
		idMap := make(map[int]struct{})
		for _, id := range course.BundleContents {
			idMap[id] = struct{}{}
		}
		courses, err := CoursesCache(GetCourses, "courses")()
		if err != nil {
			return nil, err
		}
		for _, c := range *courses {
			if _, exists := idMap[c.ID]; exists {
				fmt.Println(c.Name)
				childCourse, err := CourseCache(GetData)(c.Slug)
				if err != nil {
					return nil, err
				}
				course.Courses = append(course.Courses, childCourse)
				fmt.Println(course.Courses)
			}
		}
	}
	return course, nil
}

// GetCourses fetches the list of all courses.
func GetCourses() (*[]Course, error) {
	token, err := TokenCache(GetToken, "token")()
	if err != nil {
		return nil, fmt.Errorf("error getting token: %v", err)
	}

	url := fmt.Sprintf("https://codewithmosh.com/_next/data/%s/courses.json", token.Value)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Printf("Making GET request to: %s\n", url)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response CoursesResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &response.PageProps.Courses, nil
}

type ParentInfo struct {
	ParentName string
	Index      int
}

// GetFileNames retrieves file names for the given course.
func GetFileNames(course *Course) ([]string, error) {
	return doGetFileNames(course, nil)
}

// doGetFileNames recursively retrieves file names for the given course.
func doGetFileNames(course *Course, parent *ParentInfo) ([]string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	if course.Type == "bundle" {
		var result []string
		for i, childCourse := range course.Courses {
			filenames, err := doGetFileNames(childCourse, &ParentInfo{course.Name, i + 1})
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
				"%s/%s/%02d - %s/%02d - %s%s",
				cfg.CoursesDir,
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
