package ivmaddapi

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/internal/gen"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/stringx"
	"github.com/jzero-io/jzero/pkg/templatex"
	"github.com/spf13/cobra"
)

var (
	Name     string
	Group    string
	Handlers []string
)

func AddApi(_ *cobra.Command, _ []string) error {
	baseApiDir := filepath.Join("desc", "api")

	service := gen.GetApiServiceName(filepath.Join("desc", "api"))

	template, err := templatex.ParseTemplate(map[string]interface{}{
		"Handlers":   Handlers,
		"Service":    service,
		"Group":      Group,
		"GroupCamel": stringx.FirstUpper(stringx.ToCamel(Group)),
	}, embeded.ReadTemplateFile(filepath.Join("ivm", "add", "template.api.tpl")))
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(baseApiDir, Name+".api"), template, 0o644)
}
