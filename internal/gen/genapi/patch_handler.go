package genapi

import (
	"bytes"
	"fmt"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/jzero-io/jzero/pkg/mod"
)

type HandlerFile struct {
	Package     string
	Group       string
	Compact     bool
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
				Compact:     cast.ToBool(group.GetAnnotation("compact_handler")),
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
		file.Path = newFilePath
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

	if err = os.Rename(file.Path, newFilePath); err != nil {
		return err
	}

	file.Path = newFilePath

	// compact handler
	if file.Compact {
		err = ja.compactHandler(f, fset, file)
		if err != nil {
			return err
		}
	}

	return nil
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

func (ja *JzeroApi) compactHandler(f *ast.File, fset *token.FileSet, file HandlerFile) error {
	namingFormat, err := format.FileNamingFormat(ja.Style, filepath.Base(file.Group))
	if err != nil {
		return err
	}
	compactFile := filepath.Join(filepath.Dir(file.Path), namingFormat+".go")
	if !pathx.FileExists(compactFile) {
		_ = os.WriteFile(compactFile, []byte(fmt.Sprintf(`package %s`, f.Name.Name)), 0o644)
	}
	compactFset := token.NewFileSet()

	compactF, err := goparser.ParseFile(compactFset, compactFile, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	// 判断 compactFile 文件中是否已经存在该函数
	var (
		isExist     bool
		handlerFunc *ast.FuncDecl
	)
	for _, decl := range compactF.Decls {
		// 查找函数声明
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// 已经存在, 删掉该 handler 文件
			if funcDecl.Name.Name == util.Title(strings.TrimSuffix(file.Handler, "Handler")) {
				isExist = true
				break
			}
		}
	}
	if !isExist {
		for _, decl := range f.Decls {
			// 查找函数声明
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == util.Title(strings.TrimSuffix(file.Handler, "Handler")) {
					handlerFunc = funcDecl
				}
			}
		}
		// 将文件内容追加到compactFile
		if handlerFunc != nil {
			compactF.Decls = append(compactF.Decls, handlerFunc)
		}
		// 将 import 语句添加到 compactFile 中
		for _, i := range f.Imports {
			unquote, err := strconv.Unquote(i.Path.Value)
			if err != nil {
				return err
			}
			if i.Name != nil && i.Name.Name != "" {
				astutil.AddNamedImport(compactFset, compactF, i.Name.Name, unquote)
			} else {
				astutil.AddImport(compactFset, compactF, unquote)
			}
		}
		buf := bytes.NewBuffer(nil)
		if err := goformat.Node(buf, compactFset, compactF); err != nil {
			return err
		}
		if err = os.WriteFile(compactFile, buf.Bytes(), 0o644); err != nil {
			return err
		}
	}
	logx.Debugf("remove old handler file: %s", file.Path)
	if err = os.Remove(file.Path); err != nil {
		return err
	}

	return nil
}
