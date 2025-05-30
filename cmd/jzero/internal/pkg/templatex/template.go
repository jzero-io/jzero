package templatex

import (
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
	// "github.com/jzero-io/jzero/core/templatex"
)

var registerFuncMap = map[string]any{
	"FirstUpper": stringx.FirstUpper,
	"FirstLower": stringx.FirstLower,
	"ToCamel":    stringx.ToCamel,
}

func ParseTemplate(data any, tplT []byte) ([]byte, error) {
	// return templatex.ParseTemplate(data, tplT, templatex.WithFuncMaps([]template.FuncMap{registerFuncMap}))
	return nil, nil
}
