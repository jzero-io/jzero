package new

import (
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type JzeroApp struct {
	TemplateData map[string]interface{}
}

func (ja *JzeroApp) New() error {
	appDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app"))
	for _, file := range appDir {
		if file.IsDir() {
			continue
		}
		appFileBytes, err := templatex.ParseTemplate(ja.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", file.Name())))
		if err != nil {
			return err
		}
		appFileName := strings.TrimRight(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Output, "app", appFileName), appFileBytes)
		if err != nil {
			return err
		}
	}
	configGoFile, err := templatex.ParseTemplate(ja.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "internal", "config", "config.go.tpl")))
	if err != nil {
		return err
	}
	err = checkWrite(filepath.Join(Output, "app", "internal", "config", "config.go"), configGoFile)
	if err != nil {
		return err
	}

	middlewareDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app", "middlewares"))
	for _, file := range middlewareDir {
		if file.IsDir() {
			continue
		}
		middlewareFileBytes, err := templatex.ParseTemplate(ja.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", file.Name())))
		if err != nil {
			return err
		}
		middlewareFileName := strings.TrimRight(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Output, "app", "middlewares", middlewareFileName), middlewareFileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
