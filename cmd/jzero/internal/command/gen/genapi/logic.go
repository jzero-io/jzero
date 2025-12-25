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
	"github.com/rinchsan/gosimports"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

type LogicFile struct {
	Package      string
	Group        string
	Handler      string
	New          bool // 是否是新生成的 logic 文件
	Compact      bool // 是否合并 logic 文件
	Path         string
	DescFilepath string
	RequestType  spec.Type
	ResponseType spec.Type
}

func (ja *JzeroApi) getAllLogicFiles(apiFilepath string, apiSpec *spec.ApiSpec) ([]LogicFile, error) {
	var logicFiles []LogicFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			namingFormat, err := format.FileNamingFormat(config.C.Style, strings.TrimSuffix(route.Handler, "Handler")+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(config.C.Wd(), "internal", "logic", group.GetAnnotation("group"), namingFormat+".go")

			hf := LogicFile{
				DescFilepath: apiFilepath,
				Path:         fp,
				Group:        group.GetAnnotation("group"),
				New:          !pathx.FileExists(fp),
				Compact:      cast.ToBool(group.GetAnnotation("compact_logic")),
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

func (ja *JzeroApi) patchLogic(file LogicFile, genCodeApiSpecMap map[string]*spec.ApiSpec) error {
	// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
	newFilePath := file.Path[:len(file.Path)-8]
	// patch style
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	newFilePath += ".go"

	if pathx.FileExists(newFilePath) {
		_ = os.Remove(file.Path)
		file.Path = newFilePath
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

	if err = UpdateImportedModule(f, fset, config.C.Wd(), ja.Module); err != nil {
		return err
	}

	// change logic types
	if _, ok := genCodeApiSpecMap[file.DescFilepath]; ok {
		if err = ja.changeLogicTypes(f, fset, file); err != nil {
			console.Warning("[warning]: rewrite %s meet error %v", file.Path, err)
		}
	}

	if file.Package != "" {
		for _, g := range genCodeApiSpecMap[file.DescFilepath].Service.Groups {
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

	process, err := gosimports.Process("", buf.Bytes(), nil)
	if err != nil {
		return err
	}

	if err = os.WriteFile(file.Path, process, 0o644); err != nil {
		return err
	}

	if err = os.Rename(file.Path, newFilePath); err != nil {
		return err
	}

	file.Path = newFilePath

	return nil
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
