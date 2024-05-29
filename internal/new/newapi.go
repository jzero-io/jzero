package new

import (
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type JzeroApi struct {
	TemplateData TemplateData
	AppDir       string
}

func (ja *JzeroApi) New() error {
	apiDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app", "desc", "api"))
	for _, file := range apiDir {
		if file.IsDir() {
			continue
		}
		apiFileBytes, err := templatex.ParseTemplate(ja.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "api", file.Name())))
		if err != nil {
			return err
		}
		apiFileName := strings.TrimRight(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Output, ja.AppDir, "desc", "api", apiFileName), apiFileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
