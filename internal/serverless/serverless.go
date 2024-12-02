package serverless

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero-contrib/modx"
)

type Plugin struct {
	Path   string
	Module string
}

var PluginsTemplate = `package plugins

import (
	"github.com/zeromicro/go-zero/rest"
	
	{{range $v := .Plugins}}{{ $v.Path | base }} "{{ $v.Module }}/serverless"
	{{end}}
)

func LoadPlugins(server *rest.Server) {
	{{ range $v := .Plugins }}server.AddRoutes({{ $v.Path | base }}.Routes())
	{{end}}
}`

func GetPlugins() ([]Plugin, error) {
	wd, _ := os.Getwd()

	var plugins []Plugin
	dir, err := os.ReadDir("plugins")
	if err != nil {
		return nil, err
	}
	for _, p := range dir {
		if p.IsDir() {
			goMod, err := modx.GetGoMod(filepath.Join(wd, "plugins", p.Name()))
			if err != nil {
				return nil, err
			}
			plugins = append(plugins, Plugin{
				Path:   filepath.ToSlash(filepath.Join("plugins", p.Name())),
				Module: goMod.Path,
			})
		}
	}
	return plugins, nil
}
