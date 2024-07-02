package gen

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

	"github.com/pkg/errors"

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
	Wd                 string
	Module             string
	Style              string
	RemoveSuffix       bool
	ChangeReplaceTypes bool
}

type HandlerFile struct {
	Group   string
	Handler string
	Path    string
}

type LogicFile struct {
	Group   string
	Handler string
	Path    string

	ClientStream bool
	ServerStream bool
}

func (ja *JzeroApi) Gen() error {
	apiDirName := filepath.Join(ja.Wd, "desc", "api")

	var apiSpec *spec.ApiSpec
	// 实验性功能
	var allHandlerFiles []HandlerFile
	var allLogicFiles []LogicFile

	if !pathx.FileExists(apiDirName) {
		return nil
	}

	if pathx.FileExists(apiDirName) {
		// format api dir
		command := fmt.Sprintf("goctl api format --dir %s", apiDirName)
		_, err := execx.Run(command, ja.Wd)
		if err != nil {
			return err
		}

		fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))
		mainApiFilePath, isDelete, err := GetMainApiFilePath(apiDirName)
		if err != nil {
			return err
		}
		defer func() {
			if isDelete {
				_ = os.Remove(mainApiFilePath)
			}
		}()

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

		err = ja.generateApiCode(mainApiFilePath)
		if err != nil {
			return err
		}
		// goctl-types. make types.go separate by group
		err = ja.separateTypesGoByGoctlTypesPlugin(mainApiFilePath)
		if err != nil {
			return err
		}

		fmt.Println(color.WithColor("Done", color.FgGreen))
	}

	if ja.RemoveSuffix && apiSpec != nil {
		for _, file := range allHandlerFiles {
			if err := ja.rewriteHandlerGo(file.Path); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
			if err := ja.rewriteRoutesGo(file.Group, file.Handler); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
		}
		for _, file := range allLogicFiles {
			if err := ja.rewriteLogicGo(file.Path); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
		}
	}

	if ja.ChangeReplaceTypes {
		for _, file := range allLogicFiles {
			if err := ja.changeReplaceLogicGoTypes(file, apiSpec); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
		}
		for _, file := range allHandlerFiles {
			if err := ja.changeReplaceHandlerGoTypes(file, apiSpec); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
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
			fp := filepath.Join(ja.Wd, "internal", "handler", group.GetAnnotation("group"), namingFormat+".go")

			f := HandlerFile{
				Path:    fp,
				Group:   group.GetAnnotation("group"),
				Handler: route.Handler,
			}

			handlerFiles = append(handlerFiles, f)
		}
	}
	return handlerFiles, nil
}

func (ja *JzeroApi) getAllLogicFiles(apiSpec *spec.ApiSpec) ([]LogicFile, error) {
	var logicFiles []LogicFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			namingFormat, err := format.FileNamingFormat(ja.Style, strings.TrimSuffix(route.Handler, "Handler")+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(ja.Wd, "internal", "logic", group.GetAnnotation("group"), namingFormat+".go")

			f := LogicFile{
				Path:    fp,
				Group:   group.GetAnnotation("group"),
				Handler: route.Handler,
			}

			logicFiles = append(logicFiles, f)
		}
	}
	return logicFiles, nil
}

func getRouteApiFilePath(apiDirName string) ([]string, error) {
	var apiFiles []string

	allApiFiles, err := findApiFiles(apiDirName)
	if err != nil {
		return nil, err
	}
	for _, file := range allApiFiles {
		apiSpec, err := parser.Parse(file, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse api file %s", file)
		}
		if len(apiSpec.Service.Routes()) > 0 {
			rel, err := filepath.Rel(apiDirName, file)
			if err != nil {
				return nil, err
			}
			apiFiles = append(apiFiles, filepath.ToSlash(rel))
		}
	}

	return apiFiles, nil
}

func findApiFiles(dir string) ([]string, error) {
	var apiFiles []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			subFiles, err := findApiFiles(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			apiFiles = append(apiFiles, subFiles...)
		} else if filepath.Ext(file.Name()) == ".api" {
			apiFiles = append(apiFiles, filepath.Join(dir, file.Name()))
		}
	}

	return apiFiles, nil
}

func (ja *JzeroApi) generateApiCode(mainApiFilePath string) error {
	if mainApiFilePath == "" {
		return errors.New("empty mainApiFilePath")
	}

	fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), mainApiFilePath)
	dir := "."
	command := fmt.Sprintf("goctl api go --api %s --dir %s --home %s --style %s ", mainApiFilePath, dir, filepath.Join(embeded.Home, "go-zero"), ja.Style)
	if _, err := execx.Run(command, ja.Wd); err != nil {
		return err
	}
	return nil
}

func (ja *JzeroApi) separateTypesGoByGoctlTypesPlugin(mainApiFilePath string) error {
	dir := "."
	command := fmt.Sprintf("goctl api plugin -plugin goctl-types=\"gen\" -api %s --dir %s --style %s\n", mainApiFilePath, dir, ja.Style)
	if _, err := execx.Run(command, ja.Wd); err != nil {
		return err
	}
	return nil
}

func (ja *JzeroApi) rewriteHandlerGo(fp string) error {
	// Get the new file name of the file (without the 7 characters(Handler) before the ".go" extension)
	newFilePath := fp[:len(fp)-10]
	// patch
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	newFilePath += ".go"

	if pathx.FileExists(newFilePath) {
		_ = os.Remove(fp)
		return nil
	}

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

	return os.Rename(fp, newFilePath)
}

func (ja *JzeroApi) rewriteRoutesGo(group string, handler string) error {
	fp := filepath.Join(ja.Wd, "internal", "handler", "routes.go")
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
	// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
	newFilePath := fp[:len(fp)-8]
	// patch
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	newFilePath += ".go"

	if pathx.FileExists(newFilePath) {
		_ = os.Remove(fp)
		return nil
	}

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

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

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}

	return os.Rename(fp, newFilePath)
}

func (ja *JzeroApi) changeReplaceLogicGoTypes(file LogicFile, apiSpec *spec.ApiSpec) error {
	fp := file.Path // logic file path
	if ja.RemoveSuffix {
		// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
		fp = file.Path[:len(file.Path)-8]
		// patch
		fp = strings.TrimSuffix(fp, "_")
		fp = strings.TrimSuffix(fp, "-")
		fp += ".go"
	}

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	var methodFunc string

	var requestType string
	var responseType string

	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			if route.Handler == file.Handler && group.GetAnnotation("group") == file.Group {
				if route.RequestType != nil {
					requestType = util.Title(route.RequestType.Name())
				}
				if route.ResponseType != nil {
					responseType = util.Title(route.ResponseType.Name())
				}
				methodFunc = route.Handler
				methodFunc = util.Title(strings.TrimSuffix(methodFunc, "Handler"))
			}
		}
	}

	var needModify bool
	ast.Inspect(f, func(node ast.Node) bool {
		if fn, ok := node.(*ast.FuncDecl); ok && fn.Recv != nil {
			if fn.Name.Name == methodFunc {
				if fn.Type != nil && fn.Type.Params != nil {
					for _, param := range fn.Type.Params.List {
						if starExpr, ok := param.Type.(*ast.StarExpr); ok {
							if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
								if selectorExpr.Sel.Name != requestType {
									selectorExpr.Sel.Name = requestType
									needModify = true
								}
							}
						}
					}
				}

				if fn.Type != nil && fn.Type.Results != nil {
					for _, result := range fn.Type.Results.List {
						if starExpr, ok := result.Type.(*ast.StarExpr); ok {
							if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
								if selectorExpr.Sel.Name != responseType {
									selectorExpr.Sel.Name = responseType
									needModify = true
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	if needModify {
		// Write the modified AST back to the file
		buf := bytes.NewBuffer(nil)
		if err := goformat.Node(buf, fset, f); err != nil {
			return err
		}

		if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}

func (ja *JzeroApi) changeReplaceHandlerGoTypes(file HandlerFile, apiSpec *spec.ApiSpec) error {
	fp := file.Path // handler file path
	if ja.RemoveSuffix {
		fp = file.Path[:len(file.Path)-10]
		// patch
		fp = strings.TrimSuffix(fp, "_")
		fp = strings.TrimSuffix(fp, "-")
		fp += ".go"
	}

	var requestType string

	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			if route.Handler == file.Handler && group.GetAnnotation("group") == file.Group {
				if route.RequestType != nil {
					requestType = util.Title(route.RequestType.Name())
				}
			}
		}
	}

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	funcName := util.Title(file.Handler)
	if ja.RemoveSuffix {
		funcName = strings.TrimSuffix(funcName, "Handler")
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.Name == funcName {
				// find var req types.XXRequest
				for _, body := range fn.Body.List {
					if returnStmt, ok := body.(*ast.ReturnStmt); ok {
						for _, v := range returnStmt.Results {
							if funcLit, ok := v.(*ast.FuncLit); ok {
								for _, list := range funcLit.Body.List {
									if declStmt, ok := list.(*ast.DeclStmt); ok {
										if decl, ok := declStmt.Decl.(*ast.GenDecl); ok {
											for _, declSpec := range decl.Specs {
												if valueSpec, ok := declSpec.(*ast.ValueSpec); ok {
													if selectorExpr, ok := valueSpec.Type.(*ast.SelectorExpr); ok {
														if ident, ok := selectorExpr.X.(*ast.Ident); ok {
															if ident.Name == "types" {
																selectorExpr.Sel.Name = requestType
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
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
