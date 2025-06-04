package mod

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

// ModuleStruct contains the relative data of go module,
// which is the result of the command go list
type ModuleStruct struct {
	Path      string
	Dir       string
	GoVersion string
}

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
		if strings.Contains(string(data), "go mod tidy") {
			if _, err = execx.Run("go mod tidy", workDir); err != nil {
				return nil, err
			}
			command = exec.Command("go", "list", "-json", "-m")
			command.Dir = workDir
			if data, err = command.CombinedOutput(); err != nil {
				return nil, errors.New(string(data))
			}
		}
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
