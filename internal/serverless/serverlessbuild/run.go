package serverlessbuild

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero-contrib/templatex"
	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"

	"github.com/jzero-io/jzero/internal/serverless"
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
	pluginsGoBytes, err := templatex.ParseTemplate(map[string]any{
		"Plugins": plugins,
	}, []byte(serverless.PluginsTemplate))
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join("plugins", "plugins.go"), pluginsGoBytes, 0o644); err != nil {
		return err
	}
	return nil
}
