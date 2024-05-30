package new

import (
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type JzeroEtc struct {
	TemplateData TemplateData
	AppDir       string
}

func (je *JzeroEtc) New() error {
	cmdDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app", "etc"))
	for _, file := range cmdDir {
		if file.IsDir() {
			continue
		}
		cmdFileBytes, err := templatex.ParseTemplate(je.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "etc", file.Name())))
		if err != nil {
			return err
		}
		cmdFileName := strings.TrimSuffix(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Output, je.AppDir, "etc", cmdFileName), cmdFileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
