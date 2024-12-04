package genapi

import (
	"fmt"
	"go/ast"
	goformat "go/format"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/jzero-io/jzero/pkg/templatex"
)

func (ja *JzeroApi) separateTypesGo() error {
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
			typesGoBytes, err := templatex.ParseTemplate(map[string]any{
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

	if !ja.SplitApiTypesDir {
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
		typesGoBytes, err := templatex.ParseTemplate(map[string]any{
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

func (ja *JzeroApi) updateHandlerImportedTypesPath(f *ast.File, fset *token.FileSet, file HandlerFile) error {
	if astutil.UsesImport(f, fmt.Sprintf("%s/internal/types", ja.Module)) {
		astutil.DeleteImport(fset, f, fmt.Sprintf("%s/internal/types", ja.Module))
		astutil.AddNamedImport(fset, f, "types", fmt.Sprintf("%s/internal/types/%s", ja.Module, file.Package))
	}

	return nil
}

func (ja *JzeroApi) updateLogicImportedTypesPath(f *ast.File, fset *token.FileSet, file LogicFile) error {
	astutil.DeleteImport(fset, f, fmt.Sprintf("%s/internal/types", ja.Module))
	if file.RequestType == nil && file.ResponseType == nil {
		return nil
	}
	astutil.AddNamedImport(fset, f, "types", fmt.Sprintf("%s/internal/types/%s", ja.Module, file.Package))
	return nil
}

// changeLogicTypes just change logic file logic function params and resp, but not body and others code
func (ja *JzeroApi) changeLogicTypes(f *ast.File, fset *token.FileSet, file LogicFile) error {
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
	return nil
}
