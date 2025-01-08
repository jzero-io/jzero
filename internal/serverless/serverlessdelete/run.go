package serverlessdelete

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero-contrib/modx"
	"github.com/jzero-io/jzero-contrib/templatex"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"golang.org/x/mod/modfile"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/serverless"
	"github.com/jzero-io/jzero/pkg/mod"
)

func Run() error {
	wd, _ := os.Getwd()

	plugins, err := serverless.GetPlugins()
	if err != nil {
		return err
	}

	deletePlugins := plugins
	var remainingPlugins []serverless.Plugin

	for _, p := range config.C.Serverless.Delete.Plugin {
		deletePlugins = lo.Reject(plugins, func(item serverless.Plugin, index int) bool {
			return item.Path != filepath.ToSlash(filepath.Join("plugins", p))
		})
		remainingPlugins = lo.Filter(plugins, func(item serverless.Plugin, index int) bool {
			return item.Path != filepath.ToSlash(filepath.Join("plugins", p))
		})
	}

	if _, err := os.Stat("go.work"); err == nil {
		goWork, _ := os.ReadFile("go.work")
		work, err := modfile.ParseWork("", goWork, nil)
		if err != nil {
			return err
		}
		for _, p := range deletePlugins {
			if !strings.HasPrefix(p.Path, "./") {
				p.Path = "./" + p.Path
			}
			if err = work.DropUse(p.Path); err != nil {
				return err
			}
		}
		if err = os.WriteFile("go.work", modfile.Format(work.Syntax), 0o644); err != nil {
			return err
		}
		// reread
		goWork, _ = os.ReadFile("go.work")
		work, err = modfile.ParseWork("", goWork, nil)
		if err != nil {
			return err
		}
		if (len(work.Use) == 0) || (len(work.Use) == 1 && work.Use[0].Path == ".") {
			_ = os.Remove("go.work")
			_ = os.Remove("go.work.sum")
		}
	}

	// write plugins/plugins.go
	goMod, err := mod.GetGoMod(wd)
	if err != nil {
		return err
	}

	for i := 0; i < len(remainingPlugins); i++ {
		pluginGoMod, err := modx.GetGoMod(filepath.Join(wd, remainingPlugins[i].Path))
		if err != nil {
			return err
		}
		remainingPlugins[i].Module = pluginGoMod.Path
	}

	pluginsGoBytes, err := templatex.ParseTemplate(map[string]any{
		"Plugins": remainingPlugins,
		"Module":  goMod.Path,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("plugins", "api", "serverless_plugins.go.tpl"))))
	if err != nil {
		return err
	}
	formatBytes, err := gosimports.Process("", pluginsGoBytes, &gosimports.Options{Comments: true})
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join("plugins", "plugins.go"), formatBytes, 0o644); err != nil {
		return err
	}
	return nil
}
