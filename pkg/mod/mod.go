package mod

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/token"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/tools/go/ast/astutil"
)

// GetParentPackage if is submodule project, root package is based on go.mod and add its dir
func GetParentPackage(workDir string) (string, error) {
	mod, err := GetGoMod(workDir)
	if err != nil {
		return "", err
	}
	trim := strings.TrimPrefix(workDir, mod.Dir)
	return filepath.ToSlash(filepath.Join(mod.Path, trim)), nil
}

func GetGoVersion() (string, error) {
	resp, err := execx.Run("go env GOVERSION", "")
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(resp, "go"), nil
}

// GetGoMod is used to determine whether workDir is a go module project through command `go list -json -m`
// GetGoMod is used to determine whether workDir is a go module project through command `go list -json -m`
func GetGoMod(workDir string) (*ModuleStruct, error) {
	ms, err := GetGoMods(workDir)
	if err != nil {
		return nil, err
	}

	if len(ms) == 0 {
		return nil, errors.New("not go module project")
	}

	if len(ms) == 1 {
		return &ms[0], nil
	}

	// 是 go module 项目, 并且项目有 go.mod 文件, 但是使用了 go workspace 机制
	for _, m := range ms {
		if filepath.Clean(workDir) == filepath.Clean(m.Dir) {
			return &m, nil
		}
	}

	// 是 go module 项目. mono app 项目, 本身不存在 go.mod 文件
	// 但请保证在 go.work 中 use 中的第一行是当前项目, 如:
	// go 1.23.3
	//
	// use (
	//	.
	//	./plugins/business
	//	./plugins/resource
	// )
	return &ms[0], nil
}

func GetGoMods(workDir string) ([]ModuleStruct, error) {
	command := exec.Command("go", "list", "-json", "-m")
	command.Dir = workDir
	data, err := command.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(data))
	}

	var ms []ModuleStruct
	decoder := json.NewDecoder(bytes.NewReader(data))
	for {
		var m ModuleStruct
		err = decoder.Decode(&m)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// ModuleStruct contains the relative data of go module,
// which is the result of the command go list
type ModuleStruct struct {
	Path      string
	Dir       string
	GoVersion string
}

func UpdateImportedModule(f *ast.File, fset *token.FileSet, workDir, module string) error {
	// 当前项目存在 go.mod 项目, 并且 go list -json -m 有多个, 即使用了 go workspace 机制
	if pathx.FileExists("go.mod") {
		mods, err := GetGoMods(workDir)
		if err != nil {
			return err
		}
		if len(mods) > 1 {
			rootPkg, err := golang.GetParentPackage(workDir)
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
	return nil
}
