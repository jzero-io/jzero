package ivminit

import (
	"bytes"
	"fmt"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"strings"
)

func (ivm *IvmInit) setUpdateProtoLogic(fp string) error {
	protoParser := rpcparser.NewDefaultProtoParser()
	parse, err := protoParser.Parse(fp, true)
	if err != nil {
		return err
	}

	files, err := ivm.jzeroRpc.GetAllLogicFiles(parse)
	if err != nil {
		return err
	}

	for _, file := range files {
		newFilePath := file.Path
		if ivm.jzeroRpc.RemoveSuffix {
			// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
			newFilePath = file.Path[:len(file.Path)-8]
			// patch
			newFilePath = strings.TrimSuffix(newFilePath, "_")
			newFilePath = strings.TrimSuffix(newFilePath, "-")
			newFilePath += ".go"
		}

		fset := token.NewFileSet()

		f, err := goparser.ParseFile(fset, newFilePath, nil, goparser.ParseComments)
		if err != nil {
			return err
		}

		logicMethodName := file.Handler

		ast.Inspect(f, func(n ast.Node) bool {
			// TODO: 增加更加严格的判断
			if fn, ok := n.(*ast.FuncDecl); ok && fn.Recv != nil && fn.Name.Name == logicMethodName {
				var body []ast.Stmt
				body = append(body, &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{Name: "logic"},
					},
					Rhs: []ast.Expr{&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{Name: fmt.Sprintf("%slogic", file.Group)},
						},
					}},
				})
				body = append(body, fn.Body.List[len(fn.Body.List)-1])

				fn.Body.List = body
			}
			return true
		})
		// Write the modified AST back to the file
		buf := bytes.NewBuffer(nil)
		if err := goformat.Node(buf, fset, f); err != nil {
			return err
		}

		if err = os.WriteFile(newFilePath, buf.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
