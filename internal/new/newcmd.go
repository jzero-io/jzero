package new

import (
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type JzeroCmd struct {
	TemplateData TemplateData
}

func (jc *JzeroCmd) New() error {
	cmdDir := embeded.ReadTemplateDir(filepath.Join("jzero", "cmd"))
	for _, file := range cmdDir {
		if file.IsDir() {
			continue
		}
		cmdFileBytes, err := templatex.ParseTemplate(jc.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", file.Name())))
		if err != nil {
			return err
		}
		cmdFileName := strings.TrimRight(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Output, "cmd", cmdFileName), cmdFileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
