package plugin

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
)

type Plugin struct {
	Name   string
	Path   string
	Module string
	Mono   bool
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
				Name: p.Name(),
				Path: filepath.ToSlash(filepath.Join("plugins", p.Name())),
				Mono: func() bool {
					return !filex.FileExists(filepath.Join("plugins", p.Name(), "go.mod"))
				}(),
			})
		} else if p.Type() == os.ModeSymlink {
			plugins = append(plugins, Plugin{
				Name: p.Name(),
				Path: filepath.ToSlash(filepath.Join("plugins", p.Name())),
				Mono: func() bool {
					return !filex.FileExists(filepath.Join("plugins", p.Name(), "go.mod"))
				}(),
			})
		}
	}
	return plugins, nil
}
