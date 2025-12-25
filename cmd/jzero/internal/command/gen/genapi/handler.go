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

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/rinchsan/gosimports"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
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
			namingFormat, err := format.FileNamingFormat(config.C.Style, formatContent)
			if err != nil {
				return nil, err
			}
			fp := filepath.Join(config.C.Wd(), "internal", "handler", group.GetAnnotation("group"), namingFormat+".go")

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

func (ja *JzeroApi) patchHandler(file HandlerFile, genCodeApiSpecMap map[string]*spec.ApiSpec) error {
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

	if err = UpdateImportedModule(f, fset, config.C.Wd(), ja.Module); err != nil {
		return err
	}

	// split api types dir
	if file.Package != "" {
		for _, g := range genCodeApiSpecMap[file.ApiFilepath].Service.Groups {
			if g.GetAnnotation("group") == file.Group {
				if err = ja.updateHandlerImportedTypesPath(f, fset, file); err != nil {
					return err
				}
			}
		}
	}

	buf := bytes.NewBuffer(nil)
	if err = goformat.Node(buf, fset, f); err != nil {
		return err
	}
	process, err := gosimports.Process("", buf.Bytes(), nil)
	if err != nil {
		return err
	}

	if err := os.WriteFile(file.Path, process, 0o644); err != nil {
		return err
	}

	if err = os.Rename(file.Path, newFilePath); err != nil {
		return err
	}

	file.Path = newFilePath

	// compact handler
	df, err := decorator.ParseFile(fset, file.Path, nil, goparser.ParseComments)
	if err != nil {
		return err
	}
	if file.Compact {
		if err = ja.compactHandler(df, fset, file); err != nil {
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

func (ja *JzeroApi) compactHandler(f *dst.File, fset *token.FileSet, file HandlerFile) error {
	namingFormat, err := format.FileNamingFormat(config.C.Style, filepath.Base(file.Group))
	if err != nil {
		return err
	}
	compactFile := filepath.Join(filepath.Dir(file.Path), namingFormat+"_compact.go")
	if !pathx.FileExists(compactFile) {
		_ = os.WriteFile(compactFile, []byte(fmt.Sprintf(`package %s`, f.Name.Name)), 0o644)
	}
	compactFset := token.NewFileSet()

	compactF, err := decorator.ParseFile(compactFset, compactFile, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	// 判断 compactFile 文件中是否已经存在该函数
	var (
		isExist bool
	)
	for _, decl := range compactF.Decls {
		// 查找函数声明
		if funcDecl, ok := decl.(*dst.FuncDecl); ok {
			// 已经存在, 删掉该 handler 文件
			if funcDecl.Name.Name == util.Title(strings.TrimSuffix(file.Handler, "Handler")) {
				isExist = true
				break
			}
		}
	}
	if !isExist {
		// 将 import 语句添加到 compactFile 中
		for _, imp := range f.Imports {
			importSpec := &dst.ImportSpec{
				Path: &dst.BasicLit{
					Kind:  token.STRING,
					Value: imp.Path.Value,
				},
			}
			if imp.Name != nil {
				importSpec.Name = &dst.Ident{
					Name: imp.Name.Name,
				}
			}

			// 查找是否已经存在 import 声明
			var foundImportDecl *dst.GenDecl
			for _, decl := range compactF.Decls {
				if genDecl, ok := decl.(*dst.GenDecl); ok && genDecl.Tok == token.IMPORT {
					foundImportDecl = genDecl
					break
				}
			}

			// 如果没有找到 import 声明，创建一个新的
			if foundImportDecl == nil {
				foundImportDecl = &dst.GenDecl{
					Tok:   token.IMPORT,
					Specs: []dst.Spec{},
				}
				compactF.Decls = append([]dst.Decl{foundImportDecl}, compactF.Decls...)
			}

			// 添加导入语句
			foundImportDecl.Specs = append(foundImportDecl.Specs, importSpec)
		}
		for _, decl := range f.Decls {
			if fd, ok := decl.(*dst.FuncDecl); ok {
				compactF.Decls = append(compactF.Decls, fd)
			}

			if gd, ok := decl.(*dst.GenDecl); ok {
				if gd.Tok == token.TYPE || gd.Tok == token.VAR || gd.Tok == token.CONST {
					compactF.Decls = append(compactF.Decls, gd)
				}
			}
		}
		// 格式化并写入文件
		buf := bytes.NewBuffer(nil)
		if err := decorator.Fprint(buf, compactF); err != nil {
			return err
		}
		process, err := gosimports.Process("", buf.Bytes(), nil)
		if err != nil {
			return err
		}
		if err = os.WriteFile(compactFile, process, 0o644); err != nil {
			return err
		}
	}
	logx.Debugf("remove old handler file: %s", file.Path)
	if err = os.Remove(file.Path); err != nil {
		return err
	}

	return nil
}

func UpdateImportedModule(f *ast.File, fset *token.FileSet, workDir, module string) error {
	// 当前项目存在 go.mod 项目, 并且 go list -json -m 有多个, 即使用了 go workspace 机制
	if pathx.FileExists("go.mod") {
		mods, err := mod.GetGoMods(workDir)
		if err != nil {
			return err
		}
		if len(mods) > 1 {
			rootPkg, _, err := golang.GetParentPackage(workDir)
			if err != nil {
				return err
			}
			imports := astutil.Imports(fset, f)
			for _, imp := range imports {
				for _, name := range imp {
					if strings.HasPrefix(name.Path.Value, "\""+rootPkg) {
						unQuote, _ := strconv.Unquote(name.Path.Value)
						newImp := strings.Replace(unQuote, rootPkg, module, 1)
						astutil.RewriteImport(fset, f, unQuote, newImp)
					}
				}
			}
		}
	}
	return nil
}
