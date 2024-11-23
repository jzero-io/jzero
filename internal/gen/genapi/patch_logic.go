package genapi

import (
	"bytes"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type LogicFile struct {
	Package string

	// service
	Group string

	// rpc name
	Handler string

	Path         string
	DescFilepath string

	RequestTypeName  string
	RequestType      spec.Type
	ResponseTypeName string
	ResponseType     spec.Type
	ClientStream     bool
	ServerStream     bool
}

func (ja *JzeroApi) getAllLogicFiles(apiFilepath string, apiSpec *spec.ApiSpec) ([]LogicFile, error) {
	var logicFiles []LogicFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			namingFormat, err := format.FileNamingFormat(ja.Style, strings.TrimSuffix(route.Handler, "Handler")+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(ja.Wd, "internal", "logic", group.GetAnnotation("group"), namingFormat+".go")

			hf := LogicFile{
				DescFilepath: apiFilepath,
				Path:         fp,
				Group:        group.GetAnnotation("group"),
				Handler:      route.Handler,
				RequestType:  route.RequestType,
				ResponseType: route.ResponseType,
			}
			if goPackage, ok := apiSpec.Info.Properties["go_package"]; ok {
				hf.Package = goPackage
			}

			logicFiles = append(logicFiles, hf)
		}
	}
	return logicFiles, nil
}

func (ja *JzeroApi) patchLogic(file LogicFile) error {
	// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
	newFilePath := file.Path[:len(file.Path)-8]
	// patch style
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	newFilePath += ".go"

	if pathx.FileExists(newFilePath) {
		_ = os.Remove(file.Path)
		return nil
	}

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, file.Path, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	// remove-suffix
	if err = ja.removeLogicSuffix(f); err != nil {
		return errors.Errorf("remove suffix meet error: [%v]", err)
	}

	// change logic types
	if ja.ChangeLogicTypes {
		if _, ok := ja.GenCodeApiSpecMap[file.DescFilepath]; ok {
			if err := ja.changeLogicTypes(f, fset, file); err != nil {
				console.Warning("[warning]: rewrite %s meet error %v", file.Path, err)
			}
		}
	}

	if ja.SplitApiTypesDir {
		for _, g := range ja.GenCodeApiSpecMap[file.DescFilepath].Service.Groups {
			if g.GetAnnotation("group") == file.Group {
				if err := ja.updateLogicImportedTypesPath(f, fset, file); err != nil {
					return err
				}
			}
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err = os.WriteFile(file.Path, buf.Bytes(), 0o644); err != nil {
		return err
	}

	return os.Rename(file.Path, newFilePath)
}

func (ja *JzeroApi) removeLogicSuffix(f *ast.File) error {
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

	// change handler type struct
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
	return nil
}
