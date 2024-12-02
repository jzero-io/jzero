package serverlessdelete

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero-contrib/templatex"
	"github.com/samber/lo"
	"golang.org/x/mod/modfile"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/serverless"
	"github.com/jzero-io/jzero/internal/serverless/serverlessbuild"
)

func Run() error {
	plugins, err := serverless.GetPlugins()
	if err != nil {
		return err
	}

	for _, p := range config.C.Serverless.Delete.Plugin {
		plugins = lo.Reject(plugins, func(item serverless.Plugin, index int) bool {
			return item.Path != filepath.ToSlash(filepath.Join("plugins", p))
		})
	}

	if _, err := os.Stat("go.work"); err == nil {
		goWork, _ := os.ReadFile("go.work")
		work, err := modfile.ParseWork("", goWork, nil)
		if err != nil {
			return err
		}
		for _, p := range plugins {
			if !strings.HasPrefix(p.Module, "./") {
				p.Path = "./" + p.Path
			}
			if err = work.DropUse(p.Path); err != nil {
				return err
			}
		}
		if err = os.WriteFile("go.work", modfile.Format(work.Syntax), 0o644); err != nil {
			return err
		}
	}

	for _, p := range plugins {
		if _, err := os.Stat(p.Path); err == nil {
			if err = os.RemoveAll(p.Path); err != nil {
				return err
			}
		}
	}

	// write plugins/plugins.go
	pluginsGoBytes, err := templatex.ParseTemplate(map[string]any{
		"Plugins": plugins,
	}, []byte(serverless.PluginsTemplate))
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join("plugins", "plugins.go"), pluginsGoBytes, 0o644); err != nil {
		return err
	}
	return serverlessbuild.Run()
}
