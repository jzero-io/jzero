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
	"sync"

	"github.com/jzero-io/jzero/pkg/gitdiff"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"

	"golang.org/x/sync/errgroup"

	jgogen "github.com/jzero-io/jzero/pkg/gogen"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	zeroconfig "github.com/zeromicro/go-zero/tools/goctl/config"
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
	ApiGitDiff         bool
	ApiGitDiffPath     string
	SplitApiTypesDir   bool

	ApiFiles          []string
	GenCodeApiFiles   []string
	ApiSpecMap        map[string]*spec.ApiSpec
	GenCodeApiSpecMap map[string]*spec.ApiSpec
}

type HandlerFile struct {
	Package     string
	Group       string
	Handler     string
	Path        string
	ApiFilepath string
}

type LogicFile struct {
	Package string
	// service
	Group string
	// rpc name
	Handler     string
	Path        string
	ApiFilepath string

	RequestTypeName  string
	RequestType      spec.Type
	ResponseTypeName string
	ResponseType     spec.Type
	ClientStream     bool
	ServerStream     bool
}

func (ja *JzeroApi) Gen() error {
	apiDirName := filepath.Join("desc", "api")

	var allHandlerFiles []HandlerFile
	var allLogicFiles []LogicFile

	if !pathx.FileExists(apiDirName) {
		return nil
	}

	fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))

	// format api dir
	command := fmt.Sprintf("goctl api format --dir %s", apiDirName)
	_, err := execx.Run(command, ja.Wd)
	if err != nil {
		return err
	}

	apiFiles, err := findApiFiles(apiDirName)
	if err != nil {
		return err
	}

	ja.ApiFiles = apiFiles
	ja.ApiSpecMap = make(map[string]*spec.ApiSpec, len(apiFiles))
	ja.GenCodeApiSpecMap = make(map[string]*spec.ApiSpec, len(apiFiles))

	for _, v := range apiFiles {
		apiSpec, err := parser.Parse(v, nil)
		if err != nil {
			return err
		}
		ja.ApiSpecMap[v] = apiSpec

		logicFiles, err := ja.getAllLogicFiles(v, apiSpec)
		if err != nil {
			return err
		}
		allLogicFiles = append(allLogicFiles, logicFiles...)

		handlerFiles, err := ja.getAllHandlerFiles(v, apiSpec)
		if err != nil {
			return err
		}
		allHandlerFiles = append(allHandlerFiles, handlerFiles...)
	}

	var genCodeApiFiles []string
	if ja.ApiGitDiff {
		files, err := gitdiff.GetChangedFiles(ja.ApiGitDiffPath)
		if err == nil {
			// 获取变动的 api 文件
			genCodeApiFiles = append(genCodeApiFiles, files...)
			for _, file := range files {
				ja.GenCodeApiSpecMap[file] = ja.ApiSpecMap[file]
			}
		}
		// 获取新增的 api 文件
		files, err = gitdiff.GetAddedFiles(ja.ApiGitDiffPath)
		if err == nil {
			for _, f := range files {
				genCodeApiFiles = append(genCodeApiFiles, f)
				ja.GenCodeApiSpecMap[f] = ja.ApiSpecMap[f]
			}
		}
	} else {
		// 否则就是全量的 api 文件
		genCodeApiFiles = ja.ApiFiles
		ja.GenCodeApiSpecMap = ja.ApiSpecMap
	}
	ja.GenCodeApiFiles = genCodeApiFiles

	err = ja.generateApiCode()
	if err != nil {
		return err
	}

	// 处理多余后缀
	if ja.RemoveSuffix {
		for _, file := range allHandlerFiles {
			if _, ok := ja.GenCodeApiSpecMap[file.ApiFilepath]; ok {
				if err = ja.removeHandlerSuffix(file.Path); err != nil {
					return errors.Wrapf(err, "rewrite %s", file.Path)
				}
			}
		}
		for _, file := range allLogicFiles {
			if _, ok := ja.GenCodeApiSpecMap[file.ApiFilepath]; ok {
				if err = ja.removeLogicSuffix(file.Path); err != nil {
					return errors.Wrapf(err, "rewrite %s", file.Path)
				}
			}
		}
	}

	// 自动替换 logic 层的 request 和 response name
	if ja.ChangeReplaceTypes {
		for _, file := range allLogicFiles {
			if _, ok := ja.GenCodeApiSpecMap[file.ApiFilepath]; ok {
				if err := ja.changeLogicTypes(file); err != nil {
					console.Warning("[warning]: rewrite %s meet error %v", file.Path, err)
					continue
				}
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

func (ja *JzeroApi) generateApiCode() error {
	for _, file := range ja.GenCodeApiFiles {
		if parse, ok := ja.GenCodeApiSpecMap[file]; ok {
			for _, group := range parse.Service.Groups {
				if ja.RegenApiHandler {
					dirFile, err := os.ReadDir(filepath.Join(ja.Wd, "internal", "handler", group.GetAnnotation("group")))
					if err != nil {
						return err
					}
					for _, v := range dirFile {
						if !v.IsDir() {
							_ = os.Remove(filepath.Join(ja.Wd, "internal", "handler", group.GetAnnotation("group"), v.Name()))
						}
					}
				}
				if ja.SplitApiTypesDir {
					_ = os.RemoveAll(filepath.Join(ja.Wd, "internal", "types", group.GetAnnotation("group")))
				}
			}
		}
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

	var handlerImports ImportLines
	var allRoutesGoBody string
	var allRoutesGoBodyMap sync.Map

	var eg errgroup.Group
	eg.SetLimit(len(ja.ApiFiles))
	for _, v := range ja.ApiFiles {
		cv := v
		eg.Go(func() error {
			routesGoBody, err := ja.getRoutesGoBody(cv)
			if err != nil {
				return err
			}
			if routesGoBody != "" {
				allRoutesGoBodyMap.Store(cv, routesGoBody)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	for _, v := range ja.ApiFiles {
		if s, ok := allRoutesGoBodyMap.Load(v); ok {
			allRoutesGoBody += cast.ToString(s) + "\n"
		}
	}

	eg.SetLimit(len(ja.GenCodeApiFiles))
	for _, v := range ja.GenCodeApiFiles {
		cv := v
		if len(ja.ApiSpecMap[cv].Service.Routes()) > 0 {
			eg.Go(func() error {
				dir := "."
				fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), cv)
				// todo: 偶发的文件多线程写的 bug
				_ = os.Remove(filepath.Join("internal", "types", "types.go"))
				command := fmt.Sprintf("goctl api go --api %s --dir %s --home %s --style %s", cv, dir, goctlHome, ja.Style)
				logx.Debugf("command: %s", command)
				if _, err := execx.Run(command, ja.Wd); err != nil {
					return errors.Wrapf(err, "api file: %s", cv)
				}
				return nil
			})
		}
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	exist := make(map[string]struct{})
	for _, v := range ja.ApiFiles {
		for _, g := range ja.ApiSpecMap[v].Service.Groups {
			if _, ok := exist[g.GetAnnotation("group")]; ok {
				continue
			}
			exist[g.GetAnnotation("group")] = struct{}{}
			if g.GetAnnotation("group") != "" {
				handlerImports = append(handlerImports, fmt.Sprintf(`%s "%s/internal/handler/%s"`, strings.ToLower(strings.ReplaceAll(g.GetAnnotation("group"), "/", "")), ja.Module, g.GetAnnotation("group")))
			}
		}
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Routes":         allRoutesGoBody,
		"Module":         ja.Module,
		"HandlerImports": handlerImports,
	}, embeded.ReadTemplateFile(filepath.Join("app", "internal", "handler", "routes.go.tpl")))
	if err != nil {
		return err
	}
	source, err := goformat.Source(template)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("internal", "handler", "routes.go"), source, 0o644)
}

func (ja *JzeroApi) getRoutesGoBody(fp string) (string, error) {
	if len(ja.ApiSpecMap[fp].Service.Routes()) > 0 {
		routesGoBody, err := jgogen.GenRoutesString(ja.Module, &zeroconfig.Config{NamingFormat: ja.Style}, ja.ApiSpecMap[fp])
		if err != nil {
			return "", err
		}
		fset := token.NewFileSet()
		f, err := goparser.ParseFile(fset, "", strings.NewReader(routesGoBody), goparser.ParseComments)
		if err != nil {
			return "", err
		}
		if ja.RemoveSuffix {
			for _, g := range ja.ApiSpecMap[fp].Service.Groups {
				for _, route := range g.Routes {
					ast.Inspect(f, func(node ast.Node) bool {
						switch n := node.(type) {
						case *ast.CallExpr:
							if sel, ok := n.Fun.(*ast.SelectorExpr); ok {
								if _, ok := sel.X.(*ast.Ident); ok {
									if sel.Sel.Name == util.Title(strings.TrimSuffix(route.Handler, "Handler"))+"Handler" {
										sel.Sel.Name = util.Title(strings.TrimSuffix(route.Handler, "Handler"))
									}
								}
							} else if indent, ok := n.Fun.(*ast.Ident); ok {
								if indent.Name == util.Title(strings.TrimSuffix(route.Handler, "Handler"))+"Handler" {
									indent.Name = util.Title(strings.TrimSuffix(route.Handler, "Handler"))
								}
							}
						}
						return true
					})
				}
			}
		}
		// 遍历 AST 节点
		for _, decl := range f.Decls {
			// 查找函数声明
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == "RegisterHandlers" {
					// 提取函数体
					var buf bytes.Buffer
					if err = printer.Fprint(&buf, fset, funcDecl.Body); err != nil {
						return "", err
					}
					return buf.String(), nil
				}
			}
		}
	}
	return "", nil
}

func (ja *JzeroApi) separateTypesGo(allLogicFiles []LogicFile, allHandlerFiles []HandlerFile) error {
	_ = os.Remove(filepath.Join("internal", "types", "types.go"))

	var allTypes []spec.Type

	for _, apiFile := range ja.ApiFiles {
		allTypes = append(allTypes, ja.ApiSpecMap[apiFile].Types...)

		if ja.SplitApiTypesDir {
			typesGoString, err := gogen.BuildTypes(ja.ApiSpecMap[apiFile].Types)
			if err != nil {
				return err
			}
			goPackage, ok := ja.ApiSpecMap[apiFile].Info.Properties["go_package"]
			if !ok {
				return errors.New("do not has go_package option")
			}
			typesGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
				"Types":   typesGoString,
				"Package": strings.ToLower(strings.ReplaceAll(goPackage, "/", "")),
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
			source, err := goformat.Source(typesGoBytes)
			if err != nil {
				return err
			}
			if err = os.WriteFile(filepath.Join("internal", "types", goPackage, "types.go"), source, 0o644); err != nil {
				return err
			}
		}
	}

	if ja.SplitApiTypesDir {
		for _, v := range allLogicFiles {
			if _, ok := ja.GenCodeApiSpecMap[v.ApiFilepath]; ok {
				for _, g := range ja.GenCodeApiSpecMap[v.ApiFilepath].Service.Groups {
					if g.GetAnnotation("group") == v.Group {
						// todo 控制是否是新增的文件才更新
						if err := ja.updateLogicImportedTypesPath(v); err != nil {
							return err
						}
					}
				}
			}
		}
		for _, v := range allHandlerFiles {
			if _, ok := ja.GenCodeApiSpecMap[v.ApiFilepath]; ok {
				for _, g := range ja.GenCodeApiSpecMap[v.ApiFilepath].Service.Groups {
					if g.GetAnnotation("group") == v.Group {
						// todo 控制是否是新增的文件才更新
						if err := ja.updateHandlerImportedTypesPath(v); err != nil {
							return err
						}
					}
				}
			}
		}
	} else {
		// 去除重复
		var realAllTypes []spec.Type
		exist := make(map[string]struct{})
		for _, v := range allTypes {
			if _, ok := exist[v.Name()]; ok {
				continue
			}
			realAllTypes = append(realAllTypes, v)
			exist[v.Name()] = struct{}{}
		}

		typesGoString, err := gogen.BuildTypes(realAllTypes)
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
		source, err := goformat.Source(typesGoBytes)
		if err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join("internal", "types", "types.go"), source, 0o644); err != nil {
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
func (ja *JzeroApi) changeLogicTypes(file LogicFile) error {
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

	requestType = file.RequestType
	responseType = file.ResponseType
	methodFunc = util.Title(strings.TrimSuffix(file.Handler, "Handler"))

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
						} else if structType != nil && responseType != nil && lo.Contains(names, "w") {
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

	if astutil.UsesImport(f, fmt.Sprintf("%s/internal/types", ja.Module)) {
		astutil.DeleteImport(fset, f, fmt.Sprintf("%s/internal/types", ja.Module))
		astutil.AddNamedImport(fset, f, "types", fmt.Sprintf("%s/internal/types/%s", ja.Module, file.Package))
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

	if astutil.UsesImport(f, fmt.Sprintf("%s/internal/types", ja.Module)) {
		astutil.DeleteImport(fset, f, fmt.Sprintf("%s/internal/types", ja.Module))
		astutil.AddNamedImport(fset, f, "types", fmt.Sprintf("%s/internal/types/%s", ja.Module, file.Package))
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
