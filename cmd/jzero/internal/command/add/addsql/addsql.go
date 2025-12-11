package addsql

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func Run(args []string) error {
	baseDir := filepath.Join("desc", "sql")

	sqlName := args[0]

	template, err := templatex.ParseTemplate(filepath.Join("model", "template.sql.tpl"), map[string]any{
		"Name": sqlName,
	}, embeded.ReadTemplateFile(filepath.Join("model", "template.sql.tpl")))
	if err != nil {
		return err
	}

	_ = os.MkdirAll(filepath.Dir(filepath.Join(baseDir, sqlName)), 0755)

	if filex.FileExists(filepath.Join(baseDir, sqlName+".sql")) {
		return fmt.Errorf("%s already exists", sqlName)
	}

	err = os.WriteFile(filepath.Join(baseDir, sqlName+".sql"), template, 0o644)
	if err != nil {
		return err
	}

	return nil
}
