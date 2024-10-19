package ivmaddapi

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/format"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/stringx"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type Handler struct {
	Name string
	Verb string
}

func AddApi(c config.Config) error {
	baseApiDir := filepath.Join("desc", "api")

	service := desc.GetApiServiceName(filepath.Join("desc", "api"))

	var handlers []Handler
	for _, v := range c.Ivm.Add.Api.Handlers {
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

	template, err := templatex.ParseTemplate(map[string]any{
		"Handlers":         handlers,
		"Service":          service,
		"Group":            c.Ivm.Add.Api.Group,
		"GroupCamel":       stringx.FirstUpper(stringx.ToCamel(c.Ivm.Add.Api.Group)),
		"SplitApiTypesDir": c.Gen.SplitApiTypesDir,
	}, embeded.ReadTemplateFile(filepath.Join("ivm", "add", "template.api.tpl")))
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(baseApiDir, c.Ivm.Add.Api.Name+".api"), template, 0o644)
	if err != nil {
		return err
	}

	// format
	return format.ApiFormatByPath(filepath.Join(baseApiDir, c.Ivm.Add.Api.Name+".api"), false)
}
