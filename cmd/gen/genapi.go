package gen

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/embeded"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

type JzeroApi struct {
	Wd           string
	Module       string
	Style        string
	RemoveSuffix bool
}

type HandlerFile struct {
	Group   string
	Handler string
	Path    string
	Skip    bool
}

type LogicFile struct {
	Path string
	Skip bool
}

func (ja *JzeroApi) Gen() error {
	apiDirName := filepath.Join(ja.Wd, "app", "desc", "api")

	var apiSpec *spec.ApiSpec
	// 实验性功能
	var allHandlerFiles []HandlerFile
	var allLogicFiles []LogicFile

	if pathx.FileExists(apiDirName) {
		// format api dir
		command := fmt.Sprintf("goctl api format --dir %s", apiDirName)
		_, err := execx.Run(command, ja.Wd)
		if err != nil {
			return err
		}

		fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))
		mainApiFilePath := GetMainApiFilePath(apiDirName)
		apiSpec, err = parser.Parse(mainApiFilePath, nil)
		if err != nil {
			return err
		}

		allLogicFiles, err = ja.getAllLogicFiles(apiSpec)
		if err != nil {
			return err
		}

		allHandlerFiles, err = ja.getAllHandlerFiles(apiSpec)
		if err != nil {
			return err
		}

		err = generateApiCode(ja.Wd, mainApiFilePath, ja.Style)
		if err != nil {
			return err
		}
		// goctl-types. make types.go separate by group
		err = separateTypesGoByGoctlTypesPlugin(ja.Wd, mainApiFilePath, ja.Style)
		if err != nil {
			return err
		}
		_ = os.Remove(mainApiFilePath)
		fmt.Println(color.WithColor("Done", color.FgGreen))
	}

	if ja.RemoveSuffix && apiSpec != nil {
		for _, file := range allHandlerFiles {
			if !file.Skip {
				if err := ja.rewriteHandlerGo(file.Path); err != nil {
					return err
				}
				if err := ja.rewriteRoutesGo(file.Group, file.Handler); err != nil {
					return err
				}
			}
		}
		for _, file := range allLogicFiles {
			if !file.Skip {
				if err := ja.rewriteLogicGo(file.Path); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (ja *JzeroApi) getAllHandlerFiles(apiSpec *spec.ApiSpec) ([]HandlerFile, error) {
	var handlerFiles []HandlerFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			formatContent := strings.TrimSuffix(route.Handler, "Handler") + "Handler"
			namingFormat, err := format.FileNamingFormat(ja.Style, formatContent)
			if err != nil {
				return nil, err
			}
			fp := filepath.Join(ja.Wd, "app", "internal", "handler", group.GetAnnotation("group"), namingFormat+".go")

			f := HandlerFile{
				Path:    fp,
				Group:   group.GetAnnotation("group"),
				Handler: route.Handler,
			}

			if pathx.FileExists(fp) {
				f.Skip = true
			}

			handlerFiles = append(handlerFiles, f)
		}
	}
	return handlerFiles, nil
}

func (ja *JzeroApi) getAllLogicFiles(apiSpec *spec.ApiSpec) ([]LogicFile, error) {
	var handlerFiles []LogicFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			namingFormat, err := format.FileNamingFormat(ja.Style, strings.TrimSuffix(route.Handler, "Handler")+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(ja.Wd, "app", "internal", "logic", group.GetAnnotation("group"), namingFormat+".go")

			f := LogicFile{
				Path: fp,
			}

			if pathx.FileExists(fp) {
				f.Skip = true
			}

			handlerFiles = append(handlerFiles, f)
		}
	}
	return handlerFiles, nil
}

func getRouteApiFilePath(apiDirName string) []string {
	var apiFiles []string
	_ = filepath.Walk(apiDirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".api" {
			apiSpec, err := parser.Parse(path, nil)
			if err != nil {
				return err
			}
			if len(apiSpec.Service.Routes()) > 0 {
				rel, err := filepath.Rel(apiDirName, path)
				if err != nil {
					return err
				}
				apiFiles = append(apiFiles, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	return apiFiles
}

func generateApiCode(wd string, mainApiFilePath, style string) error {
	if mainApiFilePath == "" {
		return errors.New("empty mainApiFilePath")
	}

	fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), mainApiFilePath)
	command := fmt.Sprintf("goctl api go --api %s --dir ./app --home %s --style %s ", mainApiFilePath, filepath.Join(embeded.Home, "go-zero"), style)
	if _, err := execx.Run(command, wd); err != nil {
		return err
	}
	return nil
}

func separateTypesGoByGoctlTypesPlugin(wd string, mainApiFilePath, style string) error {
	command := fmt.Sprintf("goctl api plugin -plugin goctl-types=\"gen\" -api %s --dir ./app --style %s\n", mainApiFilePath, style)
	if _, err := execx.Run(command, wd); err != nil {
		return err
	}
	return nil
}

func (ja *JzeroApi) rewriteHandlerGo(fp string) error {
	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}
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

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err := os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}

	// Get the new file name of the file (without the 7 characters(Handler) before the ".go" extension)
	newFilePath := fp[:len(fp)-10] + ".go"

	return os.Rename(fp, newFilePath)
}

func (ja *JzeroApi) rewriteRoutesGo(group string, handler string) error {
	fp := filepath.Join(ja.Wd, "app", "internal", "handler", "routes.go")
	fset := token.NewFileSet()
	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CallExpr:
			if sel, ok := n.Fun.(*ast.SelectorExpr); ok {
				if _, ok := sel.X.(*ast.Ident); ok {
					if sel.Sel.Name == util.Title(strings.TrimSuffix(handler, "Handler"))+"Handler" {
						sel.Sel.Name = util.Title(strings.TrimSuffix(handler, "Handler"))
					}
				}
			}
		}
		return true
	})
	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err := os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}

func (ja *JzeroApi) rewriteLogicGo(fp string) error {
	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	// modify NewXXLogic
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

	// modify XXLogic Struct
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

	// modify XXLogic Struct methods receiver
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

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}

	// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
	newFilePath := fp[:len(fp)-8] + ".go"

	return os.Rename(fp, newFilePath)
}
