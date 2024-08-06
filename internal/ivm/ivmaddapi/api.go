package ivmaddapi

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/pkg/stringx"
	"github.com/jzero-io/jzero/pkg/templatex"
)

func AddApi(ic config.IvmConfig) error {
	baseApiDir := filepath.Join("desc", "api")

	service := gen.GetApiServiceName(filepath.Join("desc", "api"))

	template, err := templatex.ParseTemplate(map[string]interface{}{
		"Handlers":   ic.Add.Api.Handlers,
		"Service":    service,
		"Group":      ic.Add.Api.Group,
		"GroupCamel": stringx.FirstUpper(stringx.ToCamel(ic.Add.Api.Group)),
	}, embeded.ReadTemplateFile(filepath.Join("ivm", "add", "template.api.tpl")))
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(baseApiDir, ic.Add.Api.Name+".api"), template, 0o644)
}
