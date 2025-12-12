package utils

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/facette/natsort"
)

var videoExts = map[string]struct{}{
	".mp4":  {},
	".mov":  {},
	".avi":  {},
	".mkv":  {},
	".wmv":  {},
	".flv":  {},
	".webm": {},
}

func CopyFile(src, dst string) (int64, error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("failed to open source file: %w", err)
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			fmt.Printf("failed to close source file: %v\n", err)
		}
	}(sourceFile)

	destinationFile, err := os.Create(dst)
	if err != nil {
		return 0, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if closeErr := destinationFile.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				log.Printf("error closing destination file: %v", closeErr)
			}
		}
	}()

	bytesCopied, err := io.Copy(destinationFile, sourceFile)
	if err != nil {
		return 0, fmt.Errorf("failed to copy file contents: %w", err)
	}

	return bytesCopied, err
}

func CountVideos(source string) (int, error) {
	count := 0
	err := forEachFile(source, videoFilter, func(path string) error {
		fmt.Println(filepath.Base(path))
		count++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func MoveVideos(source string, targets []string) error {
	paths, _ := readDirNatSort(source)
	if len(paths) != len(targets) {
		return fmt.Errorf("expected %d targets, got %d", len(paths), len(targets))
	}
	for i := range len(paths) {
		err := moveFile(paths[i], targets[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func moveFile(source, dest string) error {
	dir := filepath.Dir(dest)

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	if strings.HasPrefix(source, "/tmp/") {
		_, err = CopyFile(source, dest)
		if err := os.Remove(source); err != nil {
			return err
		}
	} else {
		err = os.Rename(source, dest)
	}

	return err
}

func readDirNatSort(source string) ([]string, error) {
	paths := make([]string, 0)
	err := forEachFile(source, videoFilter, func(path string) error {
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	natsort.Sort(paths)
	return paths, nil
}

func forEachFile(
	source string,
	filter func(path string) bool,
	fn func(path string) error,
) error {
	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !filter(path) {
			return nil
		}

		return fn(path)
	})
	return err
}

func videoFilter(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := videoExts[ext]
	return ok
}
