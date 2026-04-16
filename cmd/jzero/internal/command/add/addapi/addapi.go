package addapi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/format"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func Run(args []string) (string, error) {
	baseDir := filepath.Join("desc", "api")

	apiName := args[0]

	if strings.HasSuffix(apiName, ".api") {
		apiName = strings.TrimSuffix(apiName, ".api")
	}

	target := filepath.Join(baseDir, apiName+".api")

	// fix https://github.com/jzero-io/jzero/issues/405.
	// For jzero, each api file, the server name can be different.
	template, err := templatex.ParseTemplate(filepath.Join("api", "template.api.tpl"), map[string]any{
		"Service": stringx.ToCamel(apiName),
		"Group":   apiName,
	}, embeded.ReadTemplateFile(filepath.Join("api", "template.api.tpl")))
	if err != nil {
		return target, err
	}

	if config.C.Add.Output == "file" {
		if filex.FileExists(target) {
			return target, fmt.Errorf("%s already exists", apiName)
		}

		_ = os.MkdirAll(filepath.Dir(filepath.Join(baseDir, apiName)), 0o755)

		err = os.WriteFile(target, template, 0o644)
		if err != nil {
			return target, err
		}

		// format
		return target, format.ApiFormatByPath(target, false)
	}
	fmt.Println(string(template))
	return target, nil
}
