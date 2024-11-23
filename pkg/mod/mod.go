package mod

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
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

	var ms []ModuleStruct
	decoder := json.NewDecoder(bytes.NewReader([]byte(data)))
	for {
		var m ModuleStruct
		err := decoder.Decode(&m)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}

		ms = append(ms, m)
	}

	wd, _ := os.Getwd()
	for _, m := range ms {
		if filepath.Clean(wd) == filepath.Clean(m.Dir) {
			return &m, nil
		}
	}

	// 非 go.mod 项目, 作为 sub module
	// TODO
	return &ms[0], nil
}

// ModuleStruct contains the relative data of go module,
// which is the result of the command go list
type ModuleStruct struct {
	Path      string
	Dir       string
	GoVersion string
}
