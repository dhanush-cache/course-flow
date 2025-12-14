package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// isValidPath checks for Zip Slip vulnerability
func isValidPath(fpath, dest string) bool {
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)
	return strings.HasPrefix(fpath, cleanDest)
}

// Extracts a single file from the zip archive
func extractFile(f *zip.File, dest string) error {
	fpath := filepath.Join(dest, f.Name)

	if !isValidPath(fpath, dest) {
		return fmt.Errorf("illegal file path: %s", fpath)
	}

	if f.FileInfo().IsDir() {
		return os.MkdirAll(fpath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		if err := outFile.Close(); err != nil {
			fmt.Println("zip close error:", err)
		}
	}(outFile)

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			fmt.Println("zip close error:", err)
		}
	}(rc)

	_, err = io.Copy(outFile, rc)
	return err
}

// Unzip extracts a zip archive to the specified destination folder
func Unzip(src string, dest string, onExtractFile ProgressCallback) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func(r *zip.ReadCloser) {
		if err := r.Close(); err != nil {
			fmt.Println("zip close error:", err)
		}
	}(r)

	n := len(r.File)
	for i, f := range r.File {
		err = extractFile(f, dest)
		if err != nil {
			return err
		}
		onExtractFile(i+1, n, f.Name)
	}
	return nil
}
