package new

import (
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/app/pkg/templatex"
	"github.com/jzero-io/jzero/embeded"
)

type JzeroApi struct {
	TemplateData map[string]interface{}
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
		err = checkWrite(filepath.Join(Dir, "app", "desc", "api", apiFileName), apiFileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
