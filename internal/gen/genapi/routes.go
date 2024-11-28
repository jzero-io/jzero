package genapi

import (
	"bytes"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/printer"
	"go/token"
	"path/filepath"
	"slices"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	zeroconfig "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"

	"github.com/jzero-io/jzero/embeded"
	jgogen "github.com/jzero-io/jzero/pkg/gogen"
	"github.com/jzero-io/jzero/pkg/templatex"
)

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

type Route struct {
	Group string
	spec.Route
}

func (ja *JzeroApi) genRoute2Code() ([]byte, error) {
	var routes []Route
	for _, s := range ja.ApiSpecMap {
		for _, g := range s.Service.Groups {
			for _, r := range g.Routes {
				routes = append(routes, Route{
					Group: g.GetAnnotation("group"),
					Route: r,
				})
			}
		}
	}

	slices.SortFunc(routes, func(a, b Route) int {
		if a.Path < b.Path {
			return -1
		} else if a.Path > b.Path {
			return 1
		}
		return 0
	})

	template, err := templatex.ParseTemplate(map[string]any{
		"Routes": routes,
	}, embeded.ReadTemplateFile(filepath.Join("plugins", "api", "route2code.go.tpl")))
	if err != nil {
		return nil, err
	}
	source, err := goformat.Source(template)
	if err != nil {
		return nil, err
	}
	return source, nil
}
