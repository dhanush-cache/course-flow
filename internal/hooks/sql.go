package hooks

import (
	"os"
	"path/filepath"
)

// SQL defines the hook for the 'sql' course ID.
func SQL(folderPath string) error {
	filename := "[FreeCoursesOnline Me] Code With Mosh  Complete SQL Mastery/15. Securing Databases/9- Wrap Up.mp4"
	filePath := filepath.Join(folderPath, filename)

	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}
	return nil
}
