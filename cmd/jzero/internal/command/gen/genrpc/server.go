package genrpc

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

	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type ServerFile struct {
	DescFilepath string
	Path         string
	Service      string
}

func (jr *JzeroRpc) genServer(serverImports, pbImports ImportLines, registerServers RegisterLines) error {
	if !config.C.Quiet {
		fmt.Printf("%s to generate internal/server/server.go\n", console.Green("Start"))
	}
	serverFile, err := templatex.ParseTemplate(filepath.Join("rpc", "server.go.tpl"), map[string]any{
		"Module":          jr.Module,
		"ServerImports":   serverImports,
		"PbImports":       pbImports,
		"RegisterServers": registerServers,
	}, embeded.ReadTemplateFile(filepath.Join("rpc", "server.go.tpl")))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(config.C.Wd(), "internal", "server", "server.go"), serverFile, 0o644)
	if err != nil {
		return err
	}
	if !config.C.Quiet {
		fmt.Printf("%s", console.Green("Done\n"))
	}
	return nil
}

func (jr *JzeroRpc) GetAllServerFiles(descFilepath string, protoSpec rpcparser.Proto) ([]ServerFile, error) {
	var serverFiles []ServerFile
	for _, service := range protoSpec.Service {
		namingFormat, err := format.FileNamingFormat(config.C.Style, service.Name+"Server")
		if err != nil {
			return nil, err
		}
		fp := filepath.Join(config.C.Wd(), "internal", "server", strings.ToLower(service.Name), namingFormat+".go")

		f := ServerFile{
			DescFilepath: descFilepath,
			Path:         fp,
			Service:      service.Name,
		}

		serverFiles = append(serverFiles, f)
	}
	return serverFiles, nil
}

func (jr *JzeroRpc) removeServerSuffix(fp string) error {
	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && strings.HasSuffix(fn.Name.Name, "Server") {
			fn.Name.Name = strings.TrimSuffix(fn.Name.Name, "Server")
			for _, result := range fn.Type.Results.List {
				if starExpr, ok := result.Type.(*ast.StarExpr); ok {
					if indent, ok := starExpr.X.(*ast.Ident); ok {
						indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Server"))
					}
				}
			}
			for _, body := range fn.Body.List {
				if returnStmt, ok := body.(*ast.ReturnStmt); ok {
					for _, result := range returnStmt.Results {
						if unaryExpr, ok := result.(*ast.UnaryExpr); ok {
							if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
								if indent, ok := compositeLit.Type.(*ast.Ident); ok {
									indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Server"))
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

	ast.Inspect(f, func(node ast.Node) bool {
		if fn, ok := node.(*ast.GenDecl); ok && fn.Tok == token.TYPE {
			for _, s := range fn.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					ts.Name.Name = strings.TrimSuffix(ts.Name.Name, "Server")
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
						ident.Name = util.Title(strings.TrimSuffix(ident.Name, "Server"))
					}
				}
			}
			// find handlerFunc body: l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
			for _, body := range fn.Body.List {
				if assignStmt, ok := body.(*ast.AssignStmt); ok {
					for _, rhs := range assignStmt.Rhs {
						if callExpr, ok := rhs.(*ast.CallExpr); ok {
							if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								selectorExpr.Sel.Name = strings.TrimSuffix(selectorExpr.Sel.Name, "Logic")
							}
						}
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

	// Get the new file name of the file (without the 5 characters(Server or server) before the ".go" extension)
	newFilePath := fp[:len(fp)-9]
	// patch
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	if err = os.Rename(fp, newFilePath+".go"); err != nil {
		return err
	}

	return UpdateImportedModule(newFilePath+".go", config.C.Wd(), jr.Module)
}
