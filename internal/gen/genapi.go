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

	"golang.org/x/tools/go/ast/astutil"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
)

type JzeroApi struct {
	Wd                 string
	Module             string
	Style              string
	RemoveSuffix       bool
	ChangeReplaceTypes bool
	RegenApiHandler    bool
	RegenApiTypes      bool
}

type HandlerFile struct {
	Group   string
	Handler string
	Path    string
}

type LogicFile struct {
	Package string
	// service
	Group string
	// rpc name
	Handler string
	Path    string

	RequestTypeName  string
	ResponseTypeName string
	ClientStream     bool
	ServerStream     bool
}

func (ja *JzeroApi) Gen() error {
	apiDirName := filepath.Join("desc", "api")

	var apiSpec *spec.ApiSpec
	var allHandlerFiles []HandlerFile
	var allLogicFiles []LogicFile

	if !pathx.FileExists(apiDirName) {
		return nil
	}

	var goctlHome string

	if !pathx.FileExists(filepath.Join(config.C.Gen.Home, "go-zero", "api")) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
		err = embeded.WriteTemplateDir(filepath.Join("go-zero", "api"), filepath.Join(tempDir, "api"))
		if err != nil {
			return err
		}
		goctlHome = tempDir
	} else {
		goctlHome = filepath.Join(config.C.Gen.Home, "go-zero")
	}
	logx.Debugf("goctl_home = %s", goctlHome)

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

		if ja.RegenApiHandler {
			_ = os.RemoveAll(filepath.Join(ja.Wd, "internal", "handler"))
		}

		if ja.RegenApiTypes {
			_ = os.RemoveAll(filepath.Join(ja.Wd, "internal", "types"))
		}

		err = ja.generateApiCode(mainApiFilePath, goctlHome)
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
			if err := ja.removeHandlerSuffix(file.Path); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
			if err := ja.removeRouteSuffix(file.Group, file.Handler); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
		}
		for _, file := range allLogicFiles {
			if err := ja.removeLogicSuffix(file.Path); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
		}
	}

	if ja.ChangeReplaceTypes {
		for _, file := range allLogicFiles {
			if err := ja.changeLogicTypes(file, apiSpec); err != nil {
				console.Warning("[warning]: rewrite %s meet error %v", file.Path, err)
				continue
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

func getApiFileRelPath(apiDirName string) ([]string, error) {
	var apiFiles []string

	allApiFiles, err := findApiFiles(apiDirName)
	if err != nil {
		return nil, err
	}
	for _, file := range allApiFiles {
		rel, err := filepath.Rel(apiDirName, file)
		if err != nil {
			return nil, err
		}
		apiFiles = append(apiFiles, filepath.ToSlash(rel))
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

func (ja *JzeroApi) generateApiCode(mainApiFilePath, goctlHome string) error {
	if mainApiFilePath == "" {
		return errors.New("empty mainApiFilePath")
	}

	fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), mainApiFilePath)
	dir := "."
	command := fmt.Sprintf("goctl api go --api %s --dir %s --home %s --style %s", mainApiFilePath, dir, goctlHome, ja.Style)
	logx.Debugf("command: %s", command)
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

func (ja *JzeroApi) removeHandlerSuffix(fp string) error {
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

func (ja *JzeroApi) removeRouteSuffix(group, handler string) error {
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
			} else if indent, ok := n.Fun.(*ast.Ident); ok {
				if indent.Name == util.Title(strings.TrimSuffix(handler, "Handler"))+"Handler" {
					indent.Name = util.Title(strings.TrimSuffix(handler, "Handler"))
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

func (ja *JzeroApi) removeLogicSuffix(fp string) error {
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

// changeLogicTypes just change logic file logic function params and resp, but not body and others code
func (ja *JzeroApi) changeLogicTypes(file LogicFile, apiSpec *spec.ApiSpec) error {
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

	var (
		methodFunc                string
		requestType, responseType spec.Type
	)

	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			if route.Handler == file.Handler && group.GetAnnotation("group") == file.Group {
				if route.RequestType != nil {
					requestType = route.RequestType
				}
				if route.ResponseType != nil {
					responseType = route.ResponseType
				}
				// todo: do not guess remove-suffix=false
				methodFunc = util.Title(strings.TrimSuffix(route.Handler, "Handler"))
			}
		}
	}

	ast.Inspect(f, func(node ast.Node) bool {
		if fn, ok := node.(*ast.FuncDecl); ok && fn.Recv != nil {
			if fn.Name.Name == methodFunc {
				if requestType != nil {
					switch requestType.(type) {
					case spec.DefineStruct:
						fn.Type.Params.List = []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("req")},
								Type:  &ast.StarExpr{X: ast.NewIdent("types." + requestType.Name())},
							},
						}
					}
				} else {
					fn.Type.Params.List = nil
				}

				if responseType != nil {
					switch responseType.(type) {
					case spec.PrimitiveType:
						fn.Type.Results.List = []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("resp")},
								Type:  ast.NewIdent(responseType.Name()),
							},
							{
								Names: []*ast.Ident{ast.NewIdent("err")},
								Type:  ast.NewIdent("error"),
							},
						}
					case spec.DefineStruct:
						fn.Type.Results.List = []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("resp")},
								Type:  &ast.StarExpr{X: ast.NewIdent("types." + responseType.Name())},
							},
							{
								Names: []*ast.Ident{ast.NewIdent("err")},
								Type:  ast.NewIdent("error"),
							},
						}
					}
				} else {
					fn.Type.Results.List = []*ast.Field{
						{
							Type: ast.NewIdent("error"),
						},
					}
				}
			}
		}
		return true
	})

	// change handler type struct
	var needImportNetHttp bool
	if ja.RegenApiHandler {
		ast.Inspect(f, func(node ast.Node) bool {
			if genDecl, ok := node.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, ss := range genDecl.Specs {
					if typeSpec, ok := ss.(*ast.TypeSpec); ok {
						var (
							structType *ast.StructType
							ok         bool
							names      []string
						)
						if structType, ok = typeSpec.Type.(*ast.StructType); ok {
							for _, field := range structType.Fields.List {
								for _, name := range field.Names {
									names = append(names, name.Name)
								}
							}
						}
						if structType != nil && requestType == nil && !lo.Contains(names, "r") {
							newField := &ast.Field{
								Names: []*ast.Ident{ast.NewIdent("r")},
								Type:  &ast.StarExpr{X: ast.NewIdent("http.Request")},
							}
							structType.Fields.List = append(structType.Fields.List, newField)
							needImportNetHttp = true
						} else if structType != nil && requestType != nil && lo.Contains(names, "r") {
							for i, v := range structType.Fields.List {
								if len(v.Names) > 0 {
									if v.Names[0].Name == "r" {
										// 删除这个元素
										structType.Fields.List = append(structType.Fields.List[:i], structType.Fields.List[i+1:]...)
									}
								}
							}
						}

						if structType != nil && responseType == nil && !lo.Contains(names, "w") {
							newField := &ast.Field{
								Names: []*ast.Ident{ast.NewIdent("w")},
								Type:  ast.NewIdent("http.ResponseWriter"),
							}
							structType.Fields.List = append(structType.Fields.List, newField)
							needImportNetHttp = true
						} else if structType != nil && responseType != nil && !lo.Contains(names, "w") {
							if lo.Contains(names, "w") {
								for i, v := range structType.Fields.List {
									if len(v.Names) > 0 {
										if v.Names[0].Name == "w" {
											// 删除这个元素
											structType.Fields.List = append(structType.Fields.List[:i], structType.Fields.List[i+1:]...)
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

		// change New type struct params
		ast.Inspect(f, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == fmt.Sprintf("New%s", methodFunc) {
				var paramNames []string
				for _, param := range fn.Type.Params.List {
					for _, name := range param.Names {
						paramNames = append(paramNames, name.Name)
					}
				}
				if requestType == nil && !lo.Contains(paramNames, "r") {
					fn.Type.Params.List = append(fn.Type.Params.List, &ast.Field{
						Names: []*ast.Ident{ast.NewIdent("r")},
						Type:  &ast.StarExpr{X: ast.NewIdent("http.Request")},
					})
					needImportNetHttp = true
				} else if requestType != nil && lo.Contains(paramNames, "r") {
					for i, v := range fn.Type.Params.List {
						if len(v.Names) > 0 {
							if v.Names[0].Name == "r" {
								fn.Type.Params.List = append(fn.Type.Params.List[:i], fn.Type.Params.List[i+1:]...)
							}
						}
					}
				}

				if responseType == nil && !lo.Contains(paramNames, "w") {
					fn.Type.Params.List = append(fn.Type.Params.List, &ast.Field{
						Names: []*ast.Ident{ast.NewIdent("w")},
						Type:  ast.NewIdent("http.ResponseWriter"),
					})
					needImportNetHttp = true
				} else if responseType != nil && lo.Contains(paramNames, "w") {
					for i, v := range fn.Type.Params.List {
						if len(v.Names) > 0 {
							if v.Names[0].Name == "w" {
								fn.Type.Params.List = append(fn.Type.Params.List[:i], fn.Type.Params.List[i+1:]...)
							}
						}
					}
				}

				for _, body := range fn.Body.List {
					if returnStmt, ok := body.(*ast.ReturnStmt); ok {
						for _, result := range returnStmt.Results {
							if unaryExpr, ok := result.(*ast.UnaryExpr); ok {
								if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
									if _, ok = compositeLit.Type.(*ast.Ident); ok {
										hasR := false
										hasW := false

										for _, elt := range compositeLit.Elts {
											if kv, ok := elt.(*ast.KeyValueExpr); ok {
												if key, ok := kv.Key.(*ast.Ident); ok {
													if key.Name == "r" {
														hasR = true
													}
													if key.Name == "w" {
														hasW = true
													}
												}
											}
										}

										if requestType == nil && !hasR {
											// Add new field
											newField := &ast.KeyValueExpr{
												Key:   ast.NewIdent("r"),
												Value: ast.NewIdent("r"), // or any default value you want
											}
											compositeLit.Elts = append(compositeLit.Elts, newField)
										} else if requestType != nil && hasR {
											for i, v := range compositeLit.Elts {
												if kv, ok := v.(*ast.KeyValueExpr); ok {
													if key, ok := kv.Key.(*ast.Ident); ok {
														if key.Name == "r" {
															// 删除这个元素
															compositeLit.Elts = append(compositeLit.Elts[:i], compositeLit.Elts[i+1:]...)
														}
													}
												}
											}
										}

										if responseType == nil && !hasW {
											// Add new field
											newField := &ast.KeyValueExpr{
												Key:   ast.NewIdent("w"),
												Value: ast.NewIdent("w"), // or any default value you want
											}
											compositeLit.Elts = append(compositeLit.Elts, newField)
										} else if responseType != nil && hasW {
											for i, v := range compositeLit.Elts {
												if kv, ok := v.(*ast.KeyValueExpr); ok {
													if key, ok := kv.Key.(*ast.Ident); ok {
														if key.Name == "w" {
															// 删除这个元素
															compositeLit.Elts = append(compositeLit.Elts[:i], compositeLit.Elts[i+1:]...)
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

		// check `net/http` import
		if needImportNetHttp {
			astutil.AddImport(fset, f, "net/http")
		} else {
			astutil.DeleteImport(fset, f, "net/http")
		}
	}

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}
