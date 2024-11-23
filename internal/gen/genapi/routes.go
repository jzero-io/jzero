package genapi

import (
	"bytes"
	"go/ast"
	goparser "go/parser"
	"go/printer"
	"go/token"
	"strings"

	zeroconfig "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"

	jgogen "github.com/jzero-io/jzero/pkg/gogen"
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
