package ivminit

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

	"github.com/rinchsan/gosimports"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genrpc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func (ivm *IvmInit) updateProtoLogic(fp, oldFp string) error {
	protoParser := rpcparser.NewDefaultProtoParser()
	parse, err := protoParser.Parse(fp, true)
	if err != nil {
		return err
	}

	oldParse, err := protoParser.Parse(oldFp, true)
	if err != nil {
		return err
	}

	files, err := ivm.jzeroRpc.GetAllLogicFiles(fp, parse)
	if err != nil {
		return err
	}

	oldFiles, err := ivm.jzeroRpc.GetAllLogicFiles(fp, oldParse)
	if err != nil {
		return err
	}

	for i, file := range files {
		// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
		newFilePath := file.Path[:len(file.Path)-8]
		// patch
		newFilePath = strings.TrimSuffix(newFilePath, "_")
		newFilePath = strings.TrimSuffix(newFilePath, "-")
		newFilePath += ".go"

		fset := token.NewFileSet()

		f, err := goparser.ParseFile(fset, newFilePath, nil, goparser.ParseComments)
		if err != nil {
			return err
		}

		err = ivm.astInspect(fset, f, oldFiles[i], file)
		if err != nil {
			return err
		}

		// Write the modified AST back to the file
		buf := bytes.NewBuffer(nil)
		if err := goformat.Node(buf, fset, f); err != nil {
			return err
		}

		fileContent := strings.ReplaceAll(buf.String(), "__TEMPLATE_BODY__", "{{ .Body }}")
		fileContent = strings.ReplaceAll(fileContent, "var __TEMPLATE_ADAPTOR__ string", "{{ .Adaptor }}")

		logicTypeName := file.Rpc

		templateValue := map[string]any{
			"Service":          strings.ToLower(file.Service),
			"OldService":       strings.ToLower(oldFiles[i].Service),
			"LogicTypeName":    logicTypeName,
			"MethodName":       file.Rpc,
			"RequestTypeName":  file.RequestTypeName,
			"ResponseTypeName": file.ResponseTypeName,
		}

		var templateFile []byte
		if !file.ClientStream && !file.ServerStream {
			templateLogicBody, err := templatex.ParseTemplate(templateValue, embeded.ReadTemplateFile(filepath.Join("ivm", "init", "logic-body.tpl")))
			if err != nil {
				return err
			}

			templateFile, err = templatex.ParseTemplate(map[string]any{
				"Body": string(templateLogicBody),
			}, []byte(fileContent))
			if err != nil {
				return err
			}
		} else if file.ClientStream && file.ServerStream {
			templateLogicBody, err := templatex.ParseTemplate(templateValue, embeded.ReadTemplateFile(filepath.Join("ivm", "init", "logic-client-server-stream-body.tpl")))
			if err != nil {
				return err
			}

			templateLogicAdaptor, err := templatex.ParseTemplate(templateValue, embeded.ReadTemplateFile(filepath.Join("ivm", "init", "logic-client-server-stream-adaptor.tpl")))
			if err != nil {
				return err
			}

			templateFile, err = templatex.ParseTemplate(map[string]any{
				"Body":    string(templateLogicBody),
				"Adaptor": string(templateLogicAdaptor),
			}, []byte(fileContent))
			if err != nil {
				return err
			}
		} else if file.ClientStream && !file.ServerStream {
			templateLogicBody, err := templatex.ParseTemplate(templateValue, embeded.ReadTemplateFile(filepath.Join("ivm", "init", "logic-client-stream-body.tpl")))
			if err != nil {
				return err
			}

			templateLogicAdaptor, err := templatex.ParseTemplate(templateValue, embeded.ReadTemplateFile(filepath.Join("ivm", "init", "logic-client-stream-adaptor.tpl")))
			if err != nil {
				return err
			}

			templateFile, err = templatex.ParseTemplate(map[string]any{
				"Body":    string(templateLogicBody),
				"Adaptor": string(templateLogicAdaptor),
			}, []byte(fileContent))
			if err != nil {
				return err
			}
		} else if file.ServerStream && !file.ClientStream {
			templateLogicBody, err := templatex.ParseTemplate(templateValue, embeded.ReadTemplateFile(filepath.Join("ivm", "init", "logic-server-stream-body.tpl")))
			if err != nil {
				return err
			}

			templateLogicAdaptor, err := templatex.ParseTemplate(templateValue, embeded.ReadTemplateFile(filepath.Join("ivm", "init", "logic-server-stream-adaptor.tpl")))
			if err != nil {
				return err
			}

			templateFile, err = templatex.ParseTemplate(map[string]any{
				"Body":    string(templateLogicBody),
				"Adaptor": string(templateLogicAdaptor),
			}, []byte(fileContent))
			if err != nil {
				return err
			}
		}

		templateFileFormat, err := gosimports.Process("", templateFile, &gosimports.Options{FormatOnly: true, Comments: true})
		if err != nil {
			continue
		}

		if err = os.WriteFile(newFilePath, templateFileFormat, 0o644); err != nil {
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

func (ivm *IvmInit) astAddImport(fset *token.FileSet, f *ast.File, oldService string, clientStream, serverStream bool) error {
	// 添加 import
	astutil.AddImport(fset, f, "google.golang.org/protobuf/proto")
	astutil.AddImport(fset, f, fmt.Sprintf("%s/internal/logic/%s", ivm.jzeroRpc.Module, strings.ToLower(oldService)))
	astutil.AddImport(fset, f, fmt.Sprintf("%s/internal/pb/%spb", ivm.jzeroRpc.Module, strings.ToLower(oldService)))
	if clientStream || serverStream {
		astutil.AddImport(fset, f, "io")
	}
	return nil
}

func (ivm *IvmInit) astAddLogic(fset *token.FileSet, f *ast.File, oldService, logicMethodName string, clientStream, serverStream bool) error {
	if err := ivm.astAddImport(fset, f, oldService, clientStream, serverStream); err != nil {
		return err
	}

	if clientStream || serverStream {
		varDecl := &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{
						ast.NewIdent("__TEMPLATE_ADAPTOR__"),
					},
					Type: &ast.Ident{Name: "string"},
				},
			},
		}
		f.Decls = append(f.Decls, varDecl)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Recv != nil && fn.Name.Name == logicMethodName {
			newBody := []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.Ident{
						Name: "__TEMPLATE_BODY__",
					},
				},
			}
			fn.Body.List = newBody
		}
		return true
	})

	return nil
}

func (ivm *IvmInit) astInspect(fset *token.FileSet, f *ast.File, oldLogicFile, newLogicFile genrpc.LogicFile) error {
	logicMethodName := newLogicFile.Rpc
	oldService := oldLogicFile.Service

	if err := ivm.astRemoveDefaultFirstLineComments(f); err != nil {
		return err
	}

	if err := ivm.astAddLogic(fset, f, oldService, logicMethodName, newLogicFile.ClientStream, newLogicFile.ServerStream); err != nil {
		return err
	}

	return nil
}
