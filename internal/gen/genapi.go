package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type JzeroApi struct {
	Wd                 string
	Module             string
	Style              string
	RemoveSuffix       bool
	ChangeReplaceTypes bool
	RegenApiHandler    bool
	RegenApiTypes      bool
	SplitApiTypesDir   bool
}

type HandlerFile struct {
	Package string
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

	// 处理模板
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

	// format api dir
	command := fmt.Sprintf("goctl api format --dir %s", apiDirName)
	_, err := execx.Run(command, ja.Wd)
	if err != nil {
		return err
	}

	fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))

	apiFiles, err := FindRouteApiFiles(apiDirName)
	if err != nil {
		return err
	}

	for _, v := range apiFiles {
		apiSpec, err = parser.Parse(v, nil)
		if err != nil {
			return err
		}

		logicFiles, err := ja.getAllLogicFiles(apiSpec)
		if err != nil {
			return err
		}
		allLogicFiles = append(allLogicFiles, logicFiles...)

		handlerFiles, err := ja.getAllHandlerFiles(apiSpec)
		if err != nil {
			return err
		}
		allHandlerFiles = append(allHandlerFiles, handlerFiles...)
	}

	err = ja.generateApiCode(apiFiles, goctlHome, allHandlerFiles)
	if err != nil {
		return err
	}

	// 处理多余后缀
	if ja.RemoveSuffix {
		for _, file := range allHandlerFiles {
			if err = ja.removeHandlerSuffix(file.Path); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
			if err = ja.removeRouteSuffix(file.Group, file.Handler); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
		}
		for _, file := range allLogicFiles {
			if err = ja.removeLogicSuffix(file.Path); err != nil {
				return errors.Wrapf(err, "rewrite %s", file.Path)
			}
		}
	}

	// 自动替换 logic 层的 request 和 response name
	if ja.ChangeReplaceTypes {
		for _, file := range allLogicFiles {
			if err := ja.changeLogicTypes(file, apiSpec); err != nil {
				console.Warning("[warning]: rewrite %s meet error %v", file.Path, err)
				continue
			}
		}
	}

	// 将 types.go 分 group 或者分 dir
	err = ja.separateTypesGo(allLogicFiles, allHandlerFiles)
	if err != nil {
		return err
	}

	fmt.Println(color.WithColor("Done", color.FgGreen))
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

			hf := HandlerFile{
				Path:    fp,
				Group:   group.GetAnnotation("group"),
				Handler: route.Handler,
			}
			if goPackage, ok := apiSpec.Info.Properties["go_package"]; ok {
				hf.Package = goPackage
			}
			handlerFiles = append(handlerFiles, hf)
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

			hf := LogicFile{
				Package: apiSpec.Info.Properties["go_package"],
				Path:    fp,
				Group:   group.GetAnnotation("group"),
				Handler: route.Handler,
			}
			if goPackage, ok := apiSpec.Info.Properties["go_package"]; ok {
				hf.Package = goPackage
			}

			logicFiles = append(logicFiles, hf)
		}
	}
	return logicFiles, nil
}

func (ja *JzeroApi) generateApiCode(apiFiles []string, goctlHome string, allHandlerFiles []HandlerFile) error {
	if ja.RegenApiHandler {
		_ = os.RemoveAll(filepath.Join(ja.Wd, "internal", "handler"))
	}

	if ja.RegenApiTypes {
		_ = os.RemoveAll(filepath.Join(ja.Wd, "internal", "types"))
	}

	var handlerImports ImportLines

	var allRoutesGoBody string
	for _, v := range apiFiles {
		dir := "."
		fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), v)
		command := fmt.Sprintf("goctl api go --api %s --dir %s --home %s --style %s", v, dir, goctlHome, ja.Style)
		logx.Debugf("command: %s", command)
		if _, err := execx.Run(command, ja.Wd); err != nil {
			return err
		}
		allRoutesGoBody += ja.getRoutesGoBody() + "\n"
	}

	exist := make(map[string]struct{})
	for _, v := range allHandlerFiles {
		if _, ok := exist[v.Group]; ok {
			continue
		}
		handlerImports = append(handlerImports, fmt.Sprintf(`%s "%s/internal/handler/%s"`, strings.ToLower(strings.ReplaceAll(v.Group, "/", "")), ja.Module, v.Group))
		exist[v.Group] = struct{}{}
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Routes":         allRoutesGoBody,
		"Module":         ja.Module,
		"HandlerImports": handlerImports,
	}, embeded.ReadTemplateFile(filepath.Join("app", "internal", "handler", "routes.go.tpl")))
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("internal", "handler", "routes.go"), template, 0o644)
}

func (ja *JzeroApi) getRoutesGoBody() string {
	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, filepath.Join("internal", "handler", "routes.go"), nil, goparser.ParseComments)
	if err != nil {
		return ""
	}
	// 遍历 AST 节点
	for _, decl := range f.Decls {
		// 查找函数声明
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == "RegisterHandlers" {
				// 提取函数体
				var buf bytes.Buffer
				if err = printer.Fprint(&buf, fset, funcDecl.Body); err != nil {
					return ""
				}
				return buf.String()
			}
		}
	}
	return ""
}

func (ja *JzeroApi) separateTypesGo(allLogicFiles []LogicFile, allHandlerFiles []HandlerFile) error {
	// split types go dir
	routeApiFiles, err := FindRouteApiFiles(filepath.Join("desc", "api"))
	if err != nil {
		return err
	}

	_ = os.Remove(filepath.Join("internal", "types", "types.go"))

	var allTypes []spec.Type

	for _, apiFile := range routeApiFiles {
		parse, err := parser.Parse(apiFile, "")
		if err != nil {
			return err
		}
		allTypes = append(allTypes, parse.Types...)

		if ja.SplitApiTypesDir {
			typesGoString, err := gogen.BuildTypes(parse.Types)
			if err != nil {
				return err
			}
			goPackage, ok := parse.Info.Properties["go_package"]
			if !ok {
				return errors.New("do not has go_package option")
			}
			typesGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
				"Types":   typesGoString,
				"Package": goPackage,
			}, []byte(`// Code generated by jzero. DO NOT EDIT.
package {{.Package}}

import (
    "time"
)

var (
    _ = time.Now()
)

{{.Types}}`))
			if err != nil {
				return err
			}

			_ = os.MkdirAll(filepath.Join("internal", "types", goPackage), 0o755)
			if err = os.WriteFile(filepath.Join("internal", "types", goPackage, "types.go"), typesGoBytes, 0o644); err != nil {
				return err
			}
		}
	}

	if ja.SplitApiTypesDir {
		for _, v := range allLogicFiles {
			if err = ja.updateLogicImportedTypesPath(v); err != nil {
				return err
			}
		}
		for _, v := range allHandlerFiles {
			if err = ja.updateHandlerImportedTypesPath(v); err != nil {
				return err
			}
		}
	} else {
		typesGoString, err := gogen.BuildTypes(allTypes)
		if err != nil {
			return err
		}
		typesGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
			"Types": typesGoString,
		}, []byte(`// Code generated by jzero. DO NOT EDIT.
package types

import (
    "time"
)

var (
    _ = time.Now()
)

{{.Types}}`))
		if err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join("internal", "types", "types.go"), typesGoBytes, 0o644); err != nil {
			return err
		}
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
		if requestType == nil || responseType == nil {
			astutil.AddImport(fset, f, "net/http")
		} else if requestType != nil && responseType != nil {
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

func (ja *JzeroApi) updateHandlerImportedTypesPath(file HandlerFile) error {
	fp := file.Path // handler file path
	if ja.RemoveSuffix {
		fp = file.Path[:len(file.Path)-10]
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

	astutil.DeleteImport(fset, f, fmt.Sprintf("%s/internal/types", ja.Module))
	astutil.AddNamedImport(fset, f, "types", fmt.Sprintf("%s/internal/types/%s", ja.Module, file.Group))

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

func (ja *JzeroApi) updateLogicImportedTypesPath(file LogicFile) error {
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

	astutil.DeleteImport(fset, f, fmt.Sprintf("%s/internal/types", ja.Module))
	astutil.AddNamedImport(fset, f, "types", fmt.Sprintf("%s/internal/types/%s", ja.Module, file.Group))

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
