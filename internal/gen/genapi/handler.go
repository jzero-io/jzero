package genapi

import (
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

type HandlerFile struct {
	Package     string
	Group       string
	Handler     string
	Path        string
	ApiFilepath string
}

func (ja *JzeroApi) getAllHandlerFiles(apiFilepath string, apiSpec *spec.ApiSpec) ([]HandlerFile, error) {
	var handlerFiles []HandlerFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			formatContent := strings.TrimSuffix(route.Handler, "Handler") + "Handler"
			namingFormat, err := format.FileNamingFormat(ja.Style, formatContent)
			if err != nil {
				return nil, err
			}
			fp := filepath.Join(ja.Wd, "internal", "handler", group.GetAnnotation("group"), namingFormat+".go")

			hf := HandlerFile{
				ApiFilepath: apiFilepath,
				Path:        fp,
				Group:       group.GetAnnotation("group"),
				Handler:     route.Handler,
			}
			if goPackage, ok := apiSpec.Info.Properties["go_package"]; ok {
				hf.Package = goPackage
			}
			handlerFiles = append(handlerFiles, hf)
		}
	}
	return handlerFiles, nil
}
