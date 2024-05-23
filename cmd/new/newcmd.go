package new

import (
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/app/pkg/templatex"
	"github.com/jzero-io/jzero/embeded"
)

type JzeroCmd struct {
	TemplateData map[string]interface{}
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
		err = checkWrite(filepath.Join(Dir, "cmd", cmdFileName), cmdFileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
