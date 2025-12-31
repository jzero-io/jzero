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
	"strconv"
	"strings"

	"github.com/rinchsan/gosimports"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
)

type LogicFile struct {
	Package          string
	Service          string
	Rpc              string
	Path             string
	DescFilepath     string
	RequestTypeName  string
	ResponseTypeName string
	ClientStream     bool
	ServerStream     bool
}

func (jr *JzeroRpc) GetAllLogicFiles(descFilepath string, protoSpec rpcparser.Proto) ([]LogicFile, error) {
	var logicFiles []LogicFile
	for _, service := range protoSpec.Service {
		for _, rpc := range service.RPC {
			namingFormat, err := format.FileNamingFormat(config.C.Style, rpc.Name+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(config.C.Wd(), "internal", "logic", strings.ToLower(service.Name), namingFormat+".go")

			f := LogicFile{
				Path:             fp,
				DescFilepath:     descFilepath,
				Package:          protoSpec.PbPackage,
				Rpc:              rpc.Name,
				Service:          service.Name,
				ClientStream:     rpc.StreamsRequest,
				ServerStream:     rpc.StreamsReturns,
				ResponseTypeName: rpc.ReturnsType,
				RequestTypeName:  rpc.RequestType,
			}

			logicFiles = append(logicFiles, f)
		}
	}
	return logicFiles, nil
}

func (jr *JzeroRpc) changeLogicTypes(file LogicFile) error {
	// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
	fp := file.Path[:len(file.Path)-8]
	// patch
	fp = strings.TrimSuffix(fp, "_")
	fp = strings.TrimSuffix(fp, "-")
	fp += ".go"

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(node ast.Node) bool {
		if fn, ok := node.(*ast.FuncDecl); ok && fn.Recv != nil {
			if fn.Name.Name == util.Title(file.Rpc) {
				// custom request and response
				if !file.ClientStream && !file.ServerStream {
					fn.Type.Params.List = []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("in")},
							Type:  &ast.StarExpr{X: ast.NewIdent(fmt.Sprintf("%s.%s", file.Package, util.Title(file.RequestTypeName)))},
						},
					}
					fn.Type.Results.List = []*ast.Field{
						{
							Type: &ast.StarExpr{X: ast.NewIdent(fmt.Sprintf("%s.%s", file.Package, util.Title(file.ResponseTypeName)))},
						},
						{
							Type: ast.NewIdent("error"),
						},
					}
				}

				// server stream
				if !file.ClientStream && file.ServerStream {
					fn.Type.Params.List = []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("in")},
							Type:  &ast.StarExpr{X: ast.NewIdent(fmt.Sprintf("%s.%s", file.Package, util.Title(file.RequestTypeName)))},
						},
						{
							Names: []*ast.Ident{ast.NewIdent("stream")},
							Type:  ast.NewIdent(fmt.Sprintf("%s.%s_%sServer", file.Package, util.Title(file.Service), util.Title(file.Rpc))),
						},
					}
					fn.Type.Results.List = []*ast.Field{
						{
							Type: ast.NewIdent("error"),
						},
					}
				}

				// client stream
				if file.ClientStream {
					fn.Type.Params.List = []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("stream")},
							Type:  ast.NewIdent(fmt.Sprintf("%s.%s_%sServer", file.Package, util.Title(file.Service), util.Title(file.Rpc))),
						},
					}

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

	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}
	process, err := gosimports.Process("", buf.Bytes(), nil)
	if err != nil {
		return err
	}

	if err = os.WriteFile(fp, process, 0o644); err != nil {
		return err
	}

	return nil
}

func UpdateImportedModule(filepath, workDir, module string) error {
	fset := token.NewFileSet()
	f, err := goparser.ParseFile(fset, filepath, nil, goparser.ParseComments)
	if err != nil {
		return err
	}
	// 当前项目存在 go.mod 项目, 并且 go list -json -m 有多个, 即使用了 go workspace 机制
	if pathx.FileExists("go.mod") {
		mods, err := mod.GetGoMods(workDir)
		if err != nil {
			return err
		}
		if len(mods) > 1 {
			rootPkg, _, err := golang.GetParentPackage(workDir)
			if err != nil {
				return err
			}
			imports := astutil.Imports(fset, f)
			for _, imp := range imports {
				for _, name := range imp {
					if strings.HasPrefix(name.Path.Value, "\""+rootPkg) {
						unQuote, _ := strconv.Unquote(name.Path.Value)
						newImp := strings.Replace(unQuote, rootPkg, module, 1)
						astutil.RewriteImport(fset, f, unQuote, newImp)
					}
				}
			}
		}
	}

	// write back files
	buf := bytes.NewBuffer(nil)
	if err = goformat.Node(buf, fset, f); err != nil {
		return err
	}
	process, err := gosimports.Process("", buf.Bytes(), nil)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, process, 0o644)
}

func (jr *JzeroRpc) removeLogicSuffix(fp string) error {
	// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
	newFilePath := fp[:len(fp)-8]
	// patch
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	newFilePath += ".go"

	if pathx.FileExists(newFilePath) {
		_ = os.Remove(fp)
		return nil
	}

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && strings.HasSuffix(fn.Name.Name, "Logic") {
			fn.Name.Name = strings.TrimSuffix(fn.Name.Name, "Logic")
			for _, result := range fn.Type.Results.List {
				if starExpr, ok := result.Type.(*ast.StarExpr); ok {
					if indent, ok := starExpr.X.(*ast.Ident); ok {
						indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Logic"))
					}
				}
			}
			for _, body := range fn.Body.List {
				if returnStmt, ok := body.(*ast.ReturnStmt); ok {
					for _, result := range returnStmt.Results {
						if unaryExpr, ok := result.(*ast.UnaryExpr); ok {
							if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
								if indent, ok := compositeLit.Type.(*ast.Ident); ok {
									indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Logic"))
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
					ts.Name.Name = strings.TrimSuffix(ts.Name.Name, "Logic")
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
						ident.Name = util.Title(strings.TrimSuffix(ident.Name, "Logic"))
					}
				}
			}
		}
		return true
	})

	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}
	process, err := gosimports.Process("", buf.Bytes(), nil)
	if err != nil {
		return err
	}
	if err = os.WriteFile(fp, process, 0o644); err != nil {
		return err
	}

	if err = os.Rename(fp, newFilePath); err != nil {
		return err
	}

	return UpdateImportedModule(newFilePath, config.C.Wd(), jr.Module)
}
