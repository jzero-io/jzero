package filex

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileExists check file exist
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// IsYamlFile check YAML file
func IsYamlFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

// EnsureDirExists create dir with check
func EnsureDirExists(dirPath string) error {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check directory: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", dirPath)
	}
	return nil
}
