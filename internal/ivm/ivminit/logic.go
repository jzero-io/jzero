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

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/pkg/astx"
	"github.com/jzero-io/jzero/pkg/templatex"
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

		fileContent := strings.ReplaceAll(buf.String(), "__TEMPLATE_BODY__", "{{ .Body }}")
		fileContent = strings.ReplaceAll(fileContent, "var __TEMPLATE_ADAPTOR__ string", "{{ .Adaptor }}")

		logicTypeName := fmt.Sprintf("%sLogic", file.Handler)
		if ivm.jzeroRpc.RemoveSuffix {
			logicTypeName = file.Handler
		}

		templateValue := map[string]interface{}{
			"Service":          strings.ToLower(file.Group),
			"OldService":       strings.ToLower(oldFiles[i].Group),
			"LogicTypeName":    logicTypeName,
			"MethodName":       file.Handler,
			"RequestTypeName":  file.RequestTypeName,
			"ResponseTypeName": file.ResponseTypeName,
		}

		var templateFile []byte
		if !file.ClientStream && !file.ServerStream {
			templateLogicBody, err := templatex.ParseTemplate(templateValue, embeded.ReadTemplateFile(filepath.Join("ivm", "init", "logic-body.tpl")))
			if err != nil {
				return err
			}

			templateFile, err = templatex.ParseTemplate(map[string]interface{}{
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

			templateFile, err = templatex.ParseTemplate(map[string]interface{}{
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

			templateFile, err = templatex.ParseTemplate(map[string]interface{}{
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

			templateFile, err = templatex.ParseTemplate(map[string]interface{}{
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

func (ivm *IvmInit) astAddImport(f *ast.File, oldService string, clientStream, serverStream bool) error {
	// 添加 import
	// Track added imports to avoid duplicates
	addedImports := make(map[string]bool)

	if !astx.HasImport(f, `"google.golang.org/protobuf/proto"`) {
		astx.AddImport(f, `"google.golang.org/protobuf/proto"`, addedImports)
	}

	if !astx.HasImport(f, fmt.Sprintf(`"%s/internal/logic/%s"`, ivm.jzeroRpc.Module, strings.ToLower(oldService))) {
		astx.AddImport(f, fmt.Sprintf(`"%s/internal/logic/%s"`, ivm.jzeroRpc.Module, strings.ToLower(oldService)), addedImports)
	}

	if !astx.HasImport(f, fmt.Sprintf(`"%s/internal/pb/%spb"`, ivm.jzeroRpc.Module, strings.ToLower(oldService))) {
		astx.AddImport(f, fmt.Sprintf(`"%s/internal/pb/%spb"`, ivm.jzeroRpc.Module, strings.ToLower(oldService)), addedImports)
	}
	if clientStream || serverStream {
		if !astx.HasImport(f, `"io"`) {
			astx.AddImport(f, `"io"`, addedImports)
		}
	}
	return nil
}

func (ivm *IvmInit) astAddLogic(f *ast.File, oldService, logicMethodName string, clientStream, serverStream bool) error {
	if err := ivm.astAddImport(f, oldService, clientStream, serverStream); err != nil {
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

func (ivm *IvmInit) astInspect(f *ast.File, oldLogicFile, newLogicFile gen.LogicFile) error {
	logicMethodName := newLogicFile.Handler
	oldService := oldLogicFile.Group

	if err := ivm.astRemoveDefaultFirstLineComments(f); err != nil {
		return err
	}

	if err := ivm.astAddLogic(f, oldService, logicMethodName, newLogicFile.ClientStream, newLogicFile.ServerStream); err != nil {
		return err
	}

	return nil
}
