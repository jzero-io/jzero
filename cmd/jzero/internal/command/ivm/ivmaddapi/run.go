package ivmaddapi

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type Handler struct {
	Name string
	Verb string
}

func Run() error {
	baseApiDir := filepath.Join("desc", "api")
	if !pathx.FileExists(baseApiDir) {
		_ = os.MkdirAll(baseApiDir, 0o755)
	}

	service := desc.GetApiServiceName(filepath.Join("desc", "api"))
	if service == "" {
		service = config.C.Ivm.Add.Api.Name
	}

	var handlers []Handler
	for _, v := range config.C.Ivm.Add.Api.Handlers {
		split := strings.Split(v, ":")
		var method Handler
		if len(split) == 2 {
			method.Name = split[1]
			method.Verb = split[0]
		} else if len(split) == 1 {
			method.Name = split[0]
			method.Verb = "get"
		} else {
			continue
		}
		handlers = append(handlers, method)
	}

	template, err := templatex.ParseTemplate(filepath.Join("ivm", "add", "template.api.tpl"), map[string]any{
		"Handlers":   handlers,
		"Service":    service,
		"Group":      config.C.Ivm.Add.Api.Group,
		"GroupCamel": stringx.FirstUpper(stringx.ToCamel(config.C.Ivm.Add.Api.Group)),
	}, embeded.ReadTemplateFile(filepath.Join("ivm", "add", "template.api.tpl")))
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(baseApiDir, config.C.Ivm.Add.Api.Name+".api"), template, 0o644)
	if err != nil {
		return err
	}

	// format
	return format.ApiFormatByPath(filepath.Join(baseApiDir, config.C.Ivm.Add.Api.Name+".api"), false)
}
