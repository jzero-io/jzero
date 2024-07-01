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

func (ivm *IvmInit) setUpdateProtoLogic(fp string, oldFp string) error {
	protoParser := rpcparser.NewDefaultProtoParser()
	parse, err := protoParser.Parse(fp, true)
	if err != nil {
		return err
	}

	oldParse, err := protoParser.Parse(oldFp, true)
	if err != nil {
		return err
	}

	files, err := ivm.jzeroRpc.GetAllLogicFiles(parse)
	if err != nil {
		return err
	}

	oldFiles, err := ivm.jzeroRpc.GetAllLogicFiles(oldParse)
	if err != nil {
		return err
	}

	for i, file := range files {
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
		ivm.astInspect(f, oldFiles[i].Group, file.Group, logicMethodName)

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

func (ivm *IvmInit) astInspect(f *ast.File, oldService, newService, logicMethodName string) {
	// 删除第一行注释
	if len(f.Comments) > 0 {
		// 获取第一个注释组
		firstCommentGroup := f.Comments[0]
		// 检查是否有注释
		if len(firstCommentGroup.List) > 0 {
			// 删除第一个注释
			firstCommentGroup.List = firstCommentGroup.List[1:]
			// 如果该注释组没有剩余的注释，则从文件的注释列表中删除该注释组
			if len(firstCommentGroup.List) == 0 {
				f.Comments = f.Comments[1:]
			}
		}
	}

	// 添加 import
	f.Decls = append([]ast.Decl{&ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"google.golang.org/protobuf/proto"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/logic/%s"`, ivm.jzeroRpc.Module, strings.ToLower(oldService)),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pb/%spb"`, ivm.jzeroRpc.Module, strings.ToLower(oldService)),
				},
			},
		},
	}}, f.Decls...)

	// 修改函数体逻辑
	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Recv != nil && fn.Name.Name == logicMethodName {
			// get fn request type and response type name
			var requestTypeName, responseTypeName string
			if len(fn.Type.Params.List) > 0 {
				// 第一个参数是请求类型
				requestField := fn.Type.Params.List[0]
				if field, ok := requestField.Names[0].Obj.Decl.(*ast.Field); ok {
					if startExpr, ok := field.Type.(*ast.StarExpr); ok {
						if selectorExpr, ok := startExpr.X.(*ast.SelectorExpr); ok {
							requestTypeName = selectorExpr.Sel.Name
						}
					}
				}
			}
			// 获取响应类型名称
			if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
				// 假设第一个返回值是响应类型
				responseField := fn.Type.Results.List[0]
				if starExpr, ok := responseField.Type.(*ast.StarExpr); ok {
					if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
						responseTypeName = selectorExpr.Sel.Name
					}
				}
			}

			newBody := []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "logic"}},
					Rhs: []ast.Expr{&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: fmt.Sprintf("%slogic", strings.ToLower(oldService))},
							Sel: &ast.Ident{Name: fmt.Sprintf("New%s", logicMethodName)},
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   &ast.Ident{Name: "l"},
								Sel: &ast.Ident{Name: "ctx"},
							},
							&ast.SelectorExpr{
								X:   &ast.Ident{Name: "l"},
								Sel: &ast.Ident{Name: "svcCtx"},
							},
						},
					}},
					Tok: token.DEFINE,
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "marshal"}, &ast.Ident{Name: "err"}},
					Rhs: []ast.Expr{&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "proto"},
							Sel: &ast.Ident{Name: "Marshal"},
						},
						Args: []ast.Expr{&ast.Ident{Name: "in"}},
					}},
					Tok: token.DEFINE,
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{&ast.Ident{Name: "nil"}, &ast.Ident{Name: "err"}},
							},
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{{Name: "oldIn"}},
								Type:  &ast.SelectorExpr{X: &ast.Ident{Name: fmt.Sprintf("%spb", strings.ToLower(oldService))}, Sel: &ast.Ident{Name: requestTypeName}},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "err"}},
					Rhs: []ast.Expr{&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "proto"},
							Sel: &ast.Ident{Name: "Unmarshal"},
						},
						Args: []ast.Expr{
							&ast.Ident{Name: "marshal"},
							&ast.UnaryExpr{
								Op: token.AND,
								X:  &ast.Ident{Name: "oldIn"},
							},
						},
					}},
					Tok: token.ASSIGN,
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{&ast.Ident{Name: "nil"}, &ast.Ident{Name: "err"}},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "result"}, &ast.Ident{Name: "err"}},
					Rhs: []ast.Expr{&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "logic"},
							Sel: &ast.Ident{Name: logicMethodName},
						},
						Args: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  &ast.Ident{Name: "oldIn"},
							},
						},
					}},
					Tok: token.DEFINE,
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{&ast.Ident{Name: "nil"}, &ast.Ident{Name: "err"}},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "marshal"}, &ast.Ident{Name: "err"}},
					Rhs: []ast.Expr{&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "proto"},
							Sel: &ast.Ident{Name: "Marshal"},
						},
						Args: []ast.Expr{&ast.Ident{Name: "result"}},
					}},
					Tok: token.ASSIGN,
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{&ast.Ident{Name: "nil"}, &ast.Ident{Name: "err"}},
							},
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{{Name: "newResp"}},
								Type:  &ast.SelectorExpr{X: &ast.Ident{Name: fmt.Sprintf("%spb", strings.ToLower(newService))}, Sel: &ast.Ident{Name: responseTypeName}},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "err"}},
					Rhs: []ast.Expr{&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "proto"},
							Sel: &ast.Ident{Name: "Unmarshal"},
						},
						Args: []ast.Expr{
							&ast.Ident{Name: "marshal"},
							&ast.UnaryExpr{
								Op: token.AND,
								X:  &ast.Ident{Name: "newResp"},
							},
						},
					}},
					Tok: token.ASSIGN,
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{&ast.Ident{Name: "nil"}, &ast.Ident{Name: "err"}},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X:  &ast.Ident{Name: "newResp"}},
						&ast.Ident{Name: "nil"}},
				},
			}
			fn.Body.List = newBody
		}
		return true
	})
}
