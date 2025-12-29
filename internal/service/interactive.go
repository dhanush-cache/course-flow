package service

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/db"
	"github.com/dhanush-cache/course-flow/internal/utils"
	"github.com/manifoldco/promptui"
)

func ProcessInteractive(cfg *config.Config) error {
	ctx := context.Background()

	dbDriver, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer dbDriver.Close()

	queries := db.New(dbDriver)

	course, err := askCourse(ctx, queries, cfg)
	if err != nil {
		return err
	}

	zipSource, err := askZipSource()
	if err != nil {
		return err
	}

	var zips []string
	if zipSource == ZipSourceLocal {
		zips, err = readToList("Enter paths", zips)
		if err != nil {
			return err
		}
	} else if zipSource == ZipSourceOnline {
		urls, err := queries.GetCourseURLs(ctx, course.ID)
		if err != nil {
			return err
		}

		downloadQueue := make([]string, 0)
		for _, url := range urls {
			if url.Category == "url" {
				downloadQueue = append(downloadQueue, url.Url)
			} else if url.Category == "gdrive" {
				gdriveUrl, err := utils.GDriveDirectDownloadURL(url.Url)
				if err != nil {
					return err
				}
				downloadQueue = append(downloadQueue, gdriveUrl)
			}
		}
		if len(downloadQueue) == 0 {
			downloadQueue, err = readToList("Enter urls", downloadQueue)
			if err != nil {
				return err
			}
		}

		tempDir, err := os.MkdirTemp("", "course-flow-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)

		for i, url := range downloadQueue {
			filePath := fmt.Sprintf("%s/%s", tempDir, fmt.Sprintf("file-%d.zip", i+1))
			zips = append(zips, filePath)
			err = utils.DownloadFile(url, filePath)
			if err != nil {
				return err
			}
		}
	}

	err = AddCourse(course.ID, zips, cfg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func readToList(message string, list []string) ([]string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(message + " (x for done): ")
		path, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading input: %w", err)
		}

		path = strings.TrimSpace(path)
		if path == "x" {
			break
		}

		list = append(list, path)
	}
	return list, nil
}

func askCourse(ctx context.Context, queries *db.Queries, cfg *config.Config) (*CourseWithPlatform, error) {
	courses, err := GetCoursesWithPlatform(ctx, queries, cfg)
	if err != nil {
		return nil, err
	}
	displayNames := make([]string, len(courses))
	for i, course := range courses {
		displayNames[i] = fmt.Sprintf("%s (%s)", course.Name, course.PlatformTitle)
	}

	prompt := promptui.Select{
		Label: "Select Course",
		Items: displayNames,
		Searcher: func(input string, index int) bool {
			item := strings.ToLower(displayNames[index])
			return strings.Contains(item, strings.ToLower(input))
		},
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}
	return &courses[i], nil
}

type ZipSource int

const (
	ZipSourceLocal ZipSource = iota
	ZipSourceOnline
)

func askZipSource() (ZipSource, error) {
	options := []struct {
		label  string
		source ZipSource
	}{
		{"Local File System", ZipSourceLocal},
		{"Online Download", ZipSourceOnline},
	}

	sources := make([]string, len(options))
	for i, o := range options {
		sources[i] = o.label
	}

	prompt := promptui.Select{
		Label: "Select Zip File Source",
		Items: sources,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("Prompt failed %v\n", err)
	}

	return options[i].source, nil
}
