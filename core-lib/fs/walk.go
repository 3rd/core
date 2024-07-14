package fs

import (
	"os"
	"path/filepath"
	"regexp"
)

var ignorePatterns = []string{
	`\.`,
}

func WalkFiles(path string, filter *func(path string, info os.FileInfo) bool) ([]File, error) {
	files := []File{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// skip directories
		if info.IsDir() {
			return nil
		}

		// skip ignored files
		for _, pattern := range ignorePatterns {
			if regexp.MustCompile(pattern).MatchString(path) {
				return nil
			}
		}

		// apply filter
		if filter == nil || (*filter)(path, info) {
			files = append(files, File{path, info})
		}
		return nil
	})
	return files, err
}
