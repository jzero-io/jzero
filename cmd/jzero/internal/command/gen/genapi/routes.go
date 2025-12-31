package genapi

import (
	"bytes"
	"go/ast"
	goparser "go/parser"
	"go/printer"
	"go/token"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rinchsan/gosimports"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	zeroconfig "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	jgogen "github.com/jzero-io/jzero/cmd/jzero/internal/pkg/gogen"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func (ja *JzeroApi) getRoutesGoBody(fp string, apiSpecMap map[string]*spec.ApiSpec, currentRoutesMap map[string][]spec.Route) (string, error) {
	rootPkg, projectPkg, err := golang.GetParentPackageWithModule(config.C.Wd(), ja.Module)
	if err != nil {
		return "", err
	}

	// 获取当前文件的路由（不包含 import）
	currentRoutes := currentRoutesMap[fp]
	if len(currentRoutes) == 0 {
		return "", nil
	}

	// 创建过滤后的 ApiSpec，只包含当前文件的路由
	apiSpec := apiSpecMap[fp]
	filteredSpec := &spec.ApiSpec{
		Info:    apiSpec.Info,
		Syntax:  apiSpec.Syntax,
		Imports: apiSpec.Imports,
		Types:   apiSpec.Types, // 保留所有 types（包括 import 的）
		Service: spec.Service{
			Name:   apiSpec.Service.Name,
			Groups: []spec.Group{},
		},
	}

	// 按路由路径和方法创建映射，用于快速查找
	routeKeyMap := make(map[string]bool)
	for _, r := range currentRoutes {
		key := r.Path + ":" + r.Method
		routeKeyMap[key] = true
	}

	// 过滤路由，只保留当前文件的
	for _, g := range apiSpec.Service.Groups {
		filteredGroup := spec.Group{
			Annotation: g.Annotation,
			Routes:     []spec.Route{},
		}
		for _, r := range g.Routes {
			key := r.Path + ":" + r.Method
			if routeKeyMap[key] {
				filteredGroup.Routes = append(filteredGroup.Routes, r)
			}
		}
		if len(filteredGroup.Routes) > 0 {
			filteredSpec.Service.Groups = append(filteredSpec.Service.Groups, filteredGroup)
		}
	}

	// 使用过滤后的 ApiSpec 生成路由代码
	routesGoBody, err := jgogen.GenRoutesString(rootPkg, projectPkg, &zeroconfig.Config{NamingFormat: config.C.Style}, filteredSpec)
	if err != nil {
		return "", err
	}

	fset := token.NewFileSet()
	f, err := goparser.ParseFile(fset, "", strings.NewReader(routesGoBody), goparser.ParseComments)
	if err != nil {
		return "", err
	}

	// 只处理当前文件的路由
	for _, route := range currentRoutes {
		ast.Inspect(f, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.CallExpr:
				if sel, ok := n.Fun.(*ast.SelectorExpr); ok {
					if _, ok := sel.X.(*ast.Ident); ok {
						if util.Title(sel.Sel.Name) == util.Title(strings.TrimSuffix(route.Handler, "Handler"))+"Handler" {
							sel.Sel.Name = strings.TrimSuffix(route.Handler, "Handler")
							// 查找对应的 group
							for _, g := range filteredSpec.Service.Groups {
								for _, gr := range g.Routes {
									if gr.Handler == route.Handler && g.GetAnnotation("group") != "" {
										sel.Sel.Name = util.Title(strings.TrimSuffix(route.Handler, "Handler"))
										break
									}
								}
							}
						}
					}
				} else if indent, ok := n.Fun.(*ast.Ident); ok {
					if util.Title(indent.Name) == util.Title(strings.TrimSuffix(route.Handler, "Handler"))+"Handler" {
						indent.Name = strings.TrimSuffix(route.Handler, "Handler")
						for _, g := range filteredSpec.Service.Groups {
							for _, gr := range g.Routes {
								if gr.Handler == route.Handler && g.GetAnnotation("group") != "" {
									indent.Name = util.Title(strings.TrimSuffix(route.Handler, "Handler"))
									break
								}
							}
						}
					}
				}
			}
			return true
		})
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

	return "", nil
}

type Route struct {
	Group string
	spec.Route
}

func (ja *JzeroApi) genRoute2Code(apiSpecMap map[string]*spec.ApiSpec, currentRoutesMap map[string][]spec.Route, importedFiles map[string]bool) ([]byte, error) {
	var routes []Route

	for fp, s := range apiSpecMap {
		// 跳过被 import 的文件
		if importedFiles[fp] {
			continue
		}

		// 获取当前文件的路由
		currentRoutes := currentRoutesMap[fp]
		if len(currentRoutes) == 0 {
			continue
		}

		// 为当前文件的路由创建映射
		routeKeyMap := make(map[string]bool)
		for _, r := range currentRoutes {
			key := r.Path + ":" + r.Method
			routeKeyMap[key] = true
		}

		// 只添加当前文件的路由
		for _, g := range s.Service.Groups {
			for _, r := range g.Routes {
				key := r.Path + ":" + r.Method
				if !routeKeyMap[key] {
					continue
				}

				route := Route{
					Group: g.GetAnnotation("group"),
					Route: r,
				}
				if g.GetAnnotation("prefix") != "" {
					route.Path = g.GetAnnotation("prefix") + r.Path
				}
				route.Handler = strings.TrimSuffix(r.Handler, "Handler")
				routes = append(routes, route)
			}
		}
	}

	// 先按 group 分组排序
	slices.SortFunc(routes, func(a, b Route) int {
		if a.Group < b.Group {
			return -1
		} else if a.Group > b.Group {
			return 1
		}
		return 0
	})

	// 再按 path 排序
	slices.SortStableFunc(routes, func(a, b Route) int {
		if a.Group == b.Group {
			if a.Path < b.Path {
				return -1
			} else if a.Path > b.Path {
				return 1
			}
		}
		return 0
	})

	// 最后按 method 排序
	slices.SortStableFunc(routes, func(a, b Route) int {
		if a.Group == b.Group && a.Path == b.Path {
			if a.Method < b.Method {
				return -1
			} else if a.Method > b.Method {
				return 1
			}
		}
		return 0
	})

	template, err := templatex.ParseTemplate(filepath.Join("api", "route2code.go.tpl"), map[string]any{
		"Routes": routes,
	}, embeded.ReadTemplateFile(filepath.Join("api", "route2code.go.tpl")))
	if err != nil {
		return nil, err
	}
	process, err := gosimports.Process("", template, nil)
	if err != nil {
		return nil, err
	}
	return process, nil
}
