package templatex

import (
	"text/template"

	"github.com/jzero-io/jzero/core/templatex"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
)

var registerFuncMap = map[string]any{
	"FirstUpper": stringx.FirstUpper,
	"FirstLower": stringx.FirstLower,
	"ToCamel":    stringx.ToCamel,
	"FormatStyle": func(style string, name string) string {
		namingFormat, err := format.FileNamingFormat(style, name)
		if err != nil {
			panic(err)
		}
		return namingFormat
	},
}

func ParseTemplate(data any, tplT []byte) ([]byte, error) {
	return templatex.ParseTemplate(data, tplT, templatex.WithFuncMaps([]template.FuncMap{registerFuncMap}))
}
