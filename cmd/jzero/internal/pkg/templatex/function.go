package templatex

import (
	"github.com/hashicorp/go-version"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

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
