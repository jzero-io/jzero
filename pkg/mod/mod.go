package mod

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"golang.org/x/mod/modfile"
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
func GetGoMod(workDir string) (*ModuleStruct, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}

	data, err := execx.Run("go list -json -m", workDir)
	if err != nil {
		return nil, err
	}

	var m ModuleStruct
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		// patch. 当项目存在 go.work 文件时, 为多段 json 值, 无法正常解析
		file, err := os.ReadFile(filepath.Join(workDir, "go.mod"))
		if err != nil {
			return nil, err
		}
		parse, err := modfile.Parse("", file, nil)
		if err != nil {
			return nil, err
		}
		m = ModuleStruct{
			Path:      parse.Module.Mod.Path,
			GoVersion: parse.Module.Mod.Version,
		}
		return &m, err
	}

	return &m, nil
}

// ModuleStruct contains the relative data of go module,
// which is the result of the command go list
type ModuleStruct struct {
	Path      string
	Dir       string
	GoVersion string
}
