package mod

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

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

	data, err := execx.Run("go list -json", workDir)
	if err != nil {
		return nil, err
	}

	var m ModuleStruct
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		// patch. 当项目存在 go.work 文件时, 为多段 json 值, 无法正常解析
		parse, err := modfile.Parse(filepath.Join(workDir, "go.mod"), nil, nil)
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
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}
