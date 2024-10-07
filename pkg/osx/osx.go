package osx

import "os"

func IsDir(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}
