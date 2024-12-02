package templatebuild

import (
	"os"
	"path/filepath"
)

var IgnoreDirs = []string{".git", ".idea", ".vscode", ".DS_Store", "node_modules"}

func filter(dir string, name string, ignoreDirs []string) bool {
	pwd, err := os.Getwd()
	if err != nil {
		return false
	}
	target := filepath.Join(dir, name)
	for _, id := range ignoreDirs {
		ignore := filepath.Join(pwd, id)
		if target == ignore {
			return true
		}
	}
	return false
}
