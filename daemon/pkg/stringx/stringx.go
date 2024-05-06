package stringx

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func FirstUpper(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

func GetConfigType(wd string) (string, error) {
	files, err := os.ReadDir(wd)
	if err != nil {
		return "", err
	}

	var configFile string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, _ := filepath.Match("config*", file.Name()); match {
			configFile = filepath.Join(wd, file.Name())
			break
		}
	}
	if configFile == "" {
		return "", errors.New("not found config")
	}
	return strings.TrimPrefix(filepath.Ext(configFile), "."), nil
}
