package genrpc

import (
	"bytes"
	"fmt"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/jzero-io/jzero/internal/gen/genapi"
)

func (jr *JzeroRpc) rpcStylePatchLogic(file genapi.LogicFile) error {
	fp := file.Path
	fp = fp[:len(fp)-8]
	// patch
	fp = strings.TrimSuffix(fp, "_")
	fp = strings.TrimSuffix(fp, "-")
	fp = fp + ".go"

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	packageName, _ := format.FileNamingFormat(jr.Style, file.Group)
	f.Name = ast.NewIdent(strings.ToLower(packageName))

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}

func (jr *JzeroRpc) rpcStylePatchServer(file ServerFile) error {
	fp := file.Path
	// Get the new file name of the file (without the 5 characters(Server or server) before the ".go" extension)
	fp = fp[:len(fp)-9]
	// patch
	fp = strings.TrimSuffix(fp, "_")
	fp = strings.TrimSuffix(fp, "-")
	fp = fp + ".go"

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	astutil.DeleteImport(fset, f, fmt.Sprintf("%s/internal/logic/%s", jr.Module, strings.ToLower(file.Service)))

	logicImportDir, _ := format.FileNamingFormat(jr.Style, file.Service)
	importLogicName, _ := format.FileNamingFormat("gozero", file.Service)
	astutil.AddNamedImport(fset, f, importLogicName+"logic", fmt.Sprintf("%s/internal/logic/%s", jr.Module, strings.ToLower(logicImportDir)))

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}
