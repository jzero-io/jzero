package newapifile

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/stringx"
	"github.com/jzero-io/jzero/pkg/templatex"
	"github.com/spf13/cobra"
)

var (
	Service  string
	Group    string
	Handlers []string
)

func New(_ *cobra.Command, _ []string) error {
	template, err := templatex.ParseTemplate(map[string]interface{}{
		"Handlers":   Handlers,
		"Service":    Service,
		"Group":      Group,
		"GroupCamel": stringx.FirstUpper(stringx.ToCamel(Group)),
	}, embeded.ReadTemplateFile("api-file.tpl"))
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Base(Group)+".api", template, 0o644)
}
