package genmodel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rinchsan/gosimports"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

func (js *JzeroModel) GenRegister(tables []string) error {
	var imports []string
	var tablePackages []string

	for _, t := range tables {
		imports = append(imports, fmt.Sprintf("%s/internal/model/%s", js.Module, strings.ToLower(t)))
		tablePackages = append(tablePackages, strings.ToLower(t))
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Imports":       imports,
		"TablePackages": tablePackages,
	}, embeded.ReadTemplateFile(filepath.Join("plugins", "model", "model.go.tpl")))
	if err != nil {
		return err
	}

	format, err := gosimports.Process("", template, &gosimports.Options{
		Comments:   true,
		FormatOnly: true,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("internal", "model", "model.go"), format, 0o644)
}
