package ivminit

import (
	"bytes"
	"fmt"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/jzero-io/jzero/internal/gen"

	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
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

		err = ivm.astInspect(f, oldFiles[i], file)
		if err != nil {
			return err
		}

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

func (ivm *IvmInit) astRemoveDefaultFirstLineComments(f *ast.File) error {
	if len(f.Comments) > 0 {
		firstCommentGroup := f.Comments[0]
		if len(firstCommentGroup.List) > 0 {
			firstCommentGroup.List = firstCommentGroup.List[1:]
			if len(firstCommentGroup.List) == 0 {
				f.Comments = f.Comments[1:]
			}
		}
	}
	return nil
}

func (ivm *IvmInit) astAddImport(f *ast.File, oldService string) error {
	// 添加 import
	// Track added imports to avoid duplicates
	addedImports := make(map[string]bool)

	if !hasImport(f, `"google.golang.org/protobuf/proto"`) {
		addImport(f, `"google.golang.org/protobuf/proto"`, addedImports)
	}

	if !hasImport(f, fmt.Sprintf(`"%s/internal/logic/%s"`, ivm.jzeroRpc.Module, strings.ToLower(oldService))) {
		addImport(f, fmt.Sprintf(`"%s/internal/logic/%s"`, ivm.jzeroRpc.Module, strings.ToLower(oldService)), addedImports)
	}

	if !hasImport(f, fmt.Sprintf(`"%s/internal/pb/%spb"`, ivm.jzeroRpc.Module, strings.ToLower(oldService))) {
		addImport(f, fmt.Sprintf(`"%s/internal/pb/%spb"`, ivm.jzeroRpc.Module, strings.ToLower(oldService)), addedImports)
	}
	return nil
}

func (ivm *IvmInit) astAddLogicBody(f *ast.File, oldService, newService, logicMethodName string, clientStream, serverStream bool) error {
	logicTypeName := fmt.Sprintf("New%sLogic", logicMethodName)
	if ivm.jzeroRpc.RemoveSuffix {
		logicTypeName = fmt.Sprintf("New%s", logicMethodName)
	}

	if !clientStream && !serverStream {
		ast.Inspect(f, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok && fn.Recv != nil && fn.Name.Name == logicMethodName {
				var requestTypeName, responseTypeName string
				if len(fn.Type.Params.List) > 0 {
					requestField := fn.Type.Params.List[0]
					if field, ok := requestField.Names[0].Obj.Decl.(*ast.Field); ok {
						if startExpr, ok := field.Type.(*ast.StarExpr); ok {
							if selectorExpr, ok := startExpr.X.(*ast.SelectorExpr); ok {
								requestTypeName = selectorExpr.Sel.Name
							}
						}
					}
				}
				if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
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
								Sel: &ast.Ident{Name: logicTypeName},
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
								X:  &ast.Ident{Name: "newResp"},
							},
							&ast.Ident{Name: "nil"},
						},
					},
				}
				fn.Body.List = newBody
			}
			return true
		})
	}

	return nil
}

func (ivm *IvmInit) astInspect(f *ast.File, oldLogicFile gen.LogicFile, newLogicFile gen.LogicFile) error {
	logicMethodName := newLogicFile.Handler
	oldService := oldLogicFile.Group
	newService := newLogicFile.Group

	if err := ivm.astRemoveDefaultFirstLineComments(f); err != nil {
		return err
	}

	if err := ivm.astAddImport(f, oldService); err != nil {
		return err
	}

	if err := ivm.astAddLogicBody(f, oldService, newService, logicMethodName, false, false); err != nil {
		return err
	}

	return nil
}

// hasImport checks if the given import path is already declared in the file.
func hasImport(f *ast.File, path string) bool {
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			if importSpec.Path.Value == path {
				return true
			}
		}
	}
	return false
}

// addImport adds the import declaration to the file if it's not already marked as added.
func addImport(f *ast.File, path string, addedImports map[string]bool) {
	if !addedImports[path] {
		addedImports[path] = true
		f.Decls = append([]ast.Decl{&ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: path,
					},
				},
			},
		}}, f.Decls...)
	}
}
