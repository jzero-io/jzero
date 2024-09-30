package genapi

import (
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

type LogicFile struct {
	Package string

	// service
	Group string

	// rpc name
	Handler string

	Path        string
	ApiFilepath string

	RequestTypeName  string
	RequestType      spec.Type
	ResponseTypeName string
	ResponseType     spec.Type
	ClientStream     bool
	ServerStream     bool
}

func (ja *JzeroApi) getAllLogicFiles(apiFilepath string, apiSpec *spec.ApiSpec) ([]LogicFile, error) {
	var logicFiles []LogicFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			namingFormat, err := format.FileNamingFormat(ja.Style, strings.TrimSuffix(route.Handler, "Handler")+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(ja.Wd, "internal", "logic", group.GetAnnotation("group"), namingFormat+".go")

			hf := LogicFile{
				ApiFilepath:  apiFilepath,
				Path:         fp,
				Group:        group.GetAnnotation("group"),
				Handler:      route.Handler,
				RequestType:  route.RequestType,
				ResponseType: route.ResponseType,
			}
			if goPackage, ok := apiSpec.Info.Properties["go_package"]; ok {
				hf.Package = goPackage
			}

			logicFiles = append(logicFiles, hf)
		}
	}
	return logicFiles, nil
}
