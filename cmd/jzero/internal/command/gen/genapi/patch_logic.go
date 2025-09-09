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
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

type LogicFile struct {
	Package      string
	Group        string
	Handler      string
	New          bool // 是否是新生成的 logic 文件
	Compact      bool // 是否合并 logic 文件
	Path         string
	DescFilepath string
	RequestType  spec.Type
	ResponseType spec.Type
}

func (ja *JzeroApi) getAllLogicFiles(apiFilepath string, apiSpec *spec.ApiSpec) ([]LogicFile, error) {
	var logicFiles []LogicFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			namingFormat, err := format.FileNamingFormat(config.C.Gen.Style, strings.TrimSuffix(route.Handler, "Handler")+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(config.C.Wd(), "internal", "logic", group.GetAnnotation("group"), namingFormat+".go")

			hf := LogicFile{
				DescFilepath: apiFilepath,
				Path:         fp,
				Group:        group.GetAnnotation("group"),
				New:          !pathx.FileExists(fp),
				Compact:      cast.ToBool(group.GetAnnotation("compact_logic")),
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

func (ja *JzeroApi) patchLogic(file LogicFile) error {
	// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
	newFilePath := file.Path[:len(file.Path)-8]
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
	if err = ja.removeLogicSuffix(f); err != nil {
		return errors.Errorf("remove suffix meet error: [%v]", err)
	}

	// change logic types
	if _, ok := ja.GenCodeApiSpecMap[file.DescFilepath]; ok {
		if err = ja.changeLogicTypes(f, fset, file); err != nil {
			console.Warning("[warning]: rewrite %s meet error %v", file.Path, err)
		}
	}

	if file.Package != "" {
		for _, g := range ja.GenCodeApiSpecMap[file.DescFilepath].Service.Groups {
			if g.GetAnnotation("group") == file.Group {
				if err := ja.updateLogicImportedTypesPath(f, fset, file); err != nil {
					return err
				}
			}
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	process, err := gosimports.Process("", buf.Bytes(), nil)
	if err != nil {
		return err
	}

	if err = os.WriteFile(file.Path, process, 0o644); err != nil {
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
		if err = ja.compactLogic(df, fset, file); err != nil {
			return err
		}
	}

	return nil
}

func (ja *JzeroApi) removeLogicSuffix(f *ast.File) error {
	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && strings.HasSuffix(fn.Name.Name, "Logic") {
			fn.Name.Name = strings.TrimSuffix(fn.Name.Name, "Logic")
			for _, result := range fn.Type.Results.List {
				if starExpr, ok := result.Type.(*ast.StarExpr); ok {
					if indent, ok := starExpr.X.(*ast.Ident); ok {
						indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Logic"))
					}
				}
			}
			for _, body := range fn.Body.List {
				if returnStmt, ok := body.(*ast.ReturnStmt); ok {
					for _, result := range returnStmt.Results {
						if unaryExpr, ok := result.(*ast.UnaryExpr); ok {
							if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
								if indent, ok := compositeLit.Type.(*ast.Ident); ok {
									indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Logic"))
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

	// change handler type struct
	ast.Inspect(f, func(node ast.Node) bool {
		if fn, ok := node.(*ast.GenDecl); ok && fn.Tok == token.TYPE {
			for _, s := range fn.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					ts.Name.Name = strings.TrimSuffix(ts.Name.Name, "Logic")
				}
			}
		}
		return true
	})

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Recv != nil {
			for _, list := range fn.Recv.List {
				if starExpr, ok := list.Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						ident.Name = util.Title(strings.TrimSuffix(ident.Name, "Logic"))
					}
				}
			}
		}
		return true
	})
	return nil
}

func (ja *JzeroApi) compactLogic(f *dst.File, fset *token.FileSet, file LogicFile) error {
	namingFormat, err := format.FileNamingFormat(config.C.Gen.Style, filepath.Base(file.Group))
	if err != nil {
		return err
	}
	compactFile := filepath.Join(filepath.Dir(file.Path), namingFormat+"_compact.go")
	if !pathx.FileExists(compactFile) {
		_ = os.WriteFile(compactFile, []byte(fmt.Sprintf(`package %s`, f.Name.Name)), 0o644)
	}

	// 解析目标文件
	compactF, err := decorator.ParseFile(fset, compactFile, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	// 判断 compactFile 文件中是否已经存在该函数
	var isExist bool
	for _, decl := range compactF.Decls {
		if funcDecl, ok := decl.(*dst.FuncDecl); ok {
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

		// 添加其他声明（类型、常量、变量、函数等）
		for _, decl := range f.Decls {
			if gd, ok := decl.(*dst.GenDecl); ok {
				if gd.Tok == token.TYPE || gd.Tok == token.VAR || gd.Tok == token.CONST {
					compactF.Decls = append(compactF.Decls, gd)
				}
			}
			if fd, ok := decl.(*dst.FuncDecl); ok {
				compactF.Decls = append(compactF.Decls, fd)
			}
		}

		// 格式化并写入文件
		buf := bytes.NewBuffer(nil)
		if err := decorator.Fprint(buf, compactF); err != nil {
			return err
		}
		formatted, err := goformat.Source(buf.Bytes())
		if err != nil {
			return err
		}
		process, err := gosimports.Process("", formatted, nil)
		if err != nil {
			return err
		}
		if err = os.WriteFile(compactFile, process, 0o644); err != nil {
			return err
		}
	}

	logx.Debugf("remove old logic file: %s", file.Path)
	if err = os.Remove(file.Path); err != nil {
		return err
	}

	return nil
}
