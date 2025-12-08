package utils

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// CountVideos returns the number of video files in the given directory and its subdirectories.
func CountVideos(source string) (int, error) {
	videoExts := map[string]struct{}{
		".mp4":  {},
		".mov":  {},
		".avi":  {},
		".mkv":  {},
		".wmv":  {},
		".flv":  {},
		".webm": {},
	}

	count := 0

	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if _, exists := videoExts[ext]; exists {
				count++
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
