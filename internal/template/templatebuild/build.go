package templatebuild

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

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/mod/modfile"

	"github.com/jzero-io/jzero/config"
)

func checkWrite(path string, bytes []byte) error {
	var err error
	if len(bytes) == 0 {
		return nil
	}
	if !pathx.FileExists(filepath.Dir(path)) {
		err = os.MkdirAll(filepath.Dir(path), 0o755)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(path, bytes, 0o644)
}

func Build(tc config.TemplateConfig) error {
	wd, _ := os.Getwd()

	modfileBytes, err := os.ReadFile(filepath.Join(tc.Build.WorkingDir, "go.mod"))
	if err != nil {
		return err
	}

	mod, err := modfile.Parse("", modfileBytes, nil)
	if err != nil {
		return err
	}
	tc.Build.WorkingDir = filepath.Join(wd, tc.Build.WorkingDir)
	return build(tc, tc.Build.WorkingDir, mod)
}

func build(tc config.TemplateConfig, dirname string, mod *modfile.File) error {
	dir, err := os.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, file := range dir {
		if filter(file.Name()) {
			continue
		}
		if file.IsDir() {
			err := build(tc, filepath.Join(dirname, file.Name()), mod)
			if err != nil {
				return err
			}
		} else {
			filename := fmt.Sprintf("%s.tpl", file.Name())
			fileBytes, err := os.ReadFile(filepath.Join(dirname, file.Name()))
			if err != nil {
				return err
			}
			if filepath.Ext(file.Name()) == ".go" {
				fileBytes, err = rewriteGo(mod, fileBytes)
				if err != nil {
					return err
				}
			}
			if file.Name() == "go.mod" {
				fileBytes, err = rewriteGoMod(mod, fileBytes)
				if err != nil {
					return err
				}
			}
			rel, err := filepath.Rel(tc.Build.WorkingDir, dirname)
			if err != nil {
				return err
			}
			err = checkWrite(filepath.Join(tc.Build.Output, rel, filename), fileBytes)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// rewrite go file import
func rewriteGo(mod *modfile.File, fileBytes []byte) ([]byte, error) {
	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, "", fileBytes, goparser.ParseComments)
	if err != nil {
		return nil, err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if importSpec, ok := n.(*ast.ImportSpec); ok {
			if strings.HasPrefix(importSpec.Path.Value, `"`+mod.Module.Mod.Path) {
				importSpec.Path.Value = strings.Replace(importSpec.Path.Value, mod.Module.Mod.Path, "{{ .Module }}", 1)
			}
		}
		return true
	})

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// rewrite go.mod file
func rewriteGoMod(mod *modfile.File, fileBytes []byte) ([]byte, error) {
	parse, err := modfile.Parse("", fileBytes, nil)
	if err != nil {
		return nil, err
	}

	return bytes.Replace(fileBytes, []byte(parse.Module.Mod.Path), []byte("{{ .Module }}"), 1), nil
}
