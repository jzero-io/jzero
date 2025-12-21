package templatex

import (
	"text/template"

	"github.com/hashicorp/go-version"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
)

var registerFuncMap = []template.FuncMap{
	{
		"FirstUpper":     stringx.FirstUpper,
		"FirstLower":     stringx.FirstLower,
		"ToCamel":        stringx.ToCamel,
		"FormatStyle":    FormatStyle,
		"VersionCompare": VersionCompare,
	},
}

func FormatStyle(style string, name string) string {
	namingFormat, err := format.FileNamingFormat(style, name)
	if err != nil {
		panic(err)
	}
	return namingFormat
}

func VersionCompare(v1, action, v2 string) bool {
	switch action {
	case ">":
		return version.Must(version.NewVersion(v1)).GreaterThan(version.Must(version.NewVersion(v2)))
	case "<":
		return version.Must(version.NewVersion(v1)).LessThan(version.Must(version.NewVersion(v2)))
	case ">=":
		return version.Must(version.NewVersion(v1)).GreaterThanOrEqual(version.Must(version.NewVersion(v2)))
	case "<=":
	}
	return version.Must(version.NewVersion(v1)).GreaterThanOrEqual(version.Must(version.NewVersion(v2)))
}
