package serverlessbuild

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"golang.org/x/mod/modfile"

	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/plugin"
)

func Run() error {
	wd, _ := os.Getwd()

	plugins, err := plugin.GetPlugins()
	if err != nil {
		return err
	}

	if _, err := os.Stat("go.work"); err == nil {
		goWork, _ := os.ReadFile("go.work")
		work, err := modfile.ParseWork("", goWork, nil)
		if err != nil {
			return err
		}
		for _, p := range plugins {
			if !strings.HasPrefix(p.Path, "./") {
				p.Path = "./" + p.Path
			}
			if err = work.DropUse(p.Path); err != nil {
				return err
			}
		}
		for _, p := range plugins {
			if !strings.HasPrefix(p.Path, "./") {
				p.Path = "./" + p.Path
			}
			if err = work.AddUse(p.Path, ""); err != nil {
				return err
			}
		}
		if err = os.WriteFile("go.work", modfile.Format(work.Syntax), 0o644); err != nil {
			return err
		}
	} else {
		initArgs := []string{"work", "init", "."}
		for _, p := range plugins {
			initArgs = append(initArgs, p.Path)
		}
		ec := exec.Command("go", initArgs...)
		ec.Dir = wd
		output, err := ec.CombinedOutput()
		if err != nil {
			return errors.Wrapf(err, "go work init meet error %s", string(output))
		}
	}

	// write plugins/plugins.go
	goMod, err := mod.GetGoMod(wd)
	if err != nil {
		return err
	}
	for i := 0; i < len(plugins); i++ {
		pluginGoMod, err := mod.GetGoMod(filepath.Join(wd, plugins[i].Path))
		if err != nil {
			return err
		}
		plugins[i].Module = pluginGoMod.Path
	}

	frameType, err := desc.GetFrameType()
	if err != nil {
		return err
	}

	pluginsGoBytes, err := templatex.ParseTemplate(filepath.ToSlash(filepath.Join("plugins", "api", "serverless_plugins.go.tpl")), map[string]any{
		"Plugins": plugins,
		"Module":  goMod.Path,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("plugins", frameType, "serverless_plugins.go.tpl"))))
	if err != nil {
		return err
	}
	gosimports.LocalPrefix = goMod.Path
	pluginsGoFormatBytes, err := gosimports.Process("", pluginsGoBytes, nil)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join("plugins", "plugins.go"), pluginsGoFormatBytes, 0o644); err != nil {
		return err
	}
	return nil
}
