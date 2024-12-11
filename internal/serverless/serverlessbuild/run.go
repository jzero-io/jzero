package serverlessbuild

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero-contrib/templatex"
	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"golang.org/x/mod/modfile"

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
		for _, p := range plugins {
			if !strings.HasPrefix(p.Module, "./") {
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
	pluginsGoBytes, err := templatex.ParseTemplate(map[string]any{
		"Plugins": plugins,
		"Module":  goMod.Path,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("plugins", "api", "serverless_plugins.go.tpl"))))
	if err != nil {
		return err
	}
	pluginsGoFormatBytes, err := gosimports.Process("", pluginsGoBytes, &gosimports.Options{Comments: true})
	if err := os.WriteFile(filepath.Join("plugins", "plugins.go"), pluginsGoFormatBytes, 0o644); err != nil {
		return err
	}
	return nil
}
