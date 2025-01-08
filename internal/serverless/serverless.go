package serverless

import (
	"os"
	"path/filepath"
)

type Plugin struct {
	Path   string
	Module string
}

func GetPlugins() ([]Plugin, error) {
	var plugins []Plugin
	dir, err := os.ReadDir("plugins")
	if err != nil {
		return nil, err
	}
	for _, p := range dir {
		if p.IsDir() {
			plugins = append(plugins, Plugin{
				Path: filepath.ToSlash(filepath.Join("plugins", p.Name())),
			})
		}
	}
	return plugins, nil
}
