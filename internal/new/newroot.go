package new

import (
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type JzeroRoot struct {
	TemplateData TemplateData
	AppDir       string
}

func (jr *JzeroRoot) New() error {
	appDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app"))
	for _, file := range appDir {
		if file.IsDir() {
			continue
		}
		rootFileBytes, err := templatex.ParseTemplate(jr.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", file.Name())))
		if err != nil {
			return err
		}
		rootFileName := strings.TrimSuffix(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Output, jr.AppDir, rootFileName), rootFileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
