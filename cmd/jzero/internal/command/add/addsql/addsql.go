package addsql

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func Run(args []string) error {
	baseDir := filepath.Join("desc", "sql")

	sqlName := args[0]

	if strings.HasSuffix(sqlName, ".sql") {
		sqlName = strings.TrimSuffix(sqlName, ".sql")
	}

	template, err := templatex.ParseTemplate(filepath.Join("model", "template.sql.tpl"), map[string]any{
		"Name": sqlName,
	}, embeded.ReadTemplateFile(filepath.Join("model", "template.sql.tpl")))
	if err != nil {
		return err
	}

	if config.C.Add.Output == "file" {
		if filex.FileExists(filepath.Join(baseDir, sqlName+".sql")) {
			return fmt.Errorf("%s already exists", sqlName)
		}

		_ = os.MkdirAll(filepath.Dir(filepath.Join(baseDir, sqlName)), 0o755)

		err = os.WriteFile(filepath.Join(baseDir, sqlName+".sql"), template, 0o644)
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println(string(template))
	return nil
}
