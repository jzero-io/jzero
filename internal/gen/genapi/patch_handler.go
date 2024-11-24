package genapi

import (
	"bytes"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/pkg/mod"
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

func (ja *JzeroApi) patchHandler(file HandlerFile) error {
	// Get the new file name of the file (without the 7 characters(Handler) before the ".go" extension)
	newFilePath := file.Path[:len(file.Path)-10]
	// patch style
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	newFilePath += ".go"

	if pathx.FileExists(newFilePath) {
		_ = os.Remove(file.Path)
		return nil
	}

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, file.Path, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	// remove-suffix
	if err = ja.removeHandlerSuffix(f); err != nil {
		return err
	}

	if err = mod.UpdateImportedModule(f, fset, ja.Wd, ja.Module); err != nil {
		return err
	}

	// split api types dir
	if ja.SplitApiTypesDir {
		for _, g := range ja.GenCodeApiSpecMap[file.ApiFilepath].Service.Groups {
			if g.GetAnnotation("group") == file.Group {
				if err = ja.updateHandlerImportedTypesPath(f, fset, file); err != nil {
					return err
				}
			}
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err := os.WriteFile(file.Path, buf.Bytes(), 0o644); err != nil {
		return err
	}

	return os.Rename(file.Path, newFilePath)
}

func (ja *JzeroApi) removeHandlerSuffix(f *ast.File) error {
	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && strings.HasSuffix(fn.Name.Name, "Handler") {
			fn.Name.Name = strings.TrimSuffix(fn.Name.Name, "Handler")

			// find handlerFunc body: l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
			for _, body := range fn.Body.List {
				if returnStmt, ok := body.(*ast.ReturnStmt); ok {
					for _, v := range returnStmt.Results {
						if funcLit, ok := v.(*ast.FuncLit); ok {
							for _, list := range funcLit.Body.List {
								if assignStmt, ok := list.(*ast.AssignStmt); ok {
									for _, rh := range assignStmt.Rhs {
										if callExpr, ok := rh.(*ast.CallExpr); ok {
											if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
												selectorExpr.Sel.Name = strings.TrimSuffix(selectorExpr.Sel.Name, "Logic")
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return false
		}
		return true
	})
	return nil
}
