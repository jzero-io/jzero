package mod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/mod/modfile"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
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
	// 判断是否有 go.mod, 如果有直接获取 module
	if pathx.FileExists(filepath.Join(workDir, "go.mod")) {
		// 解析 go.mod 获取 module 信息
		goModBytes, err := os.ReadFile(filepath.Join(workDir, "go.mod"))
		if err != nil {
			return nil, err
		}
		mod, err := modfile.Parse("", goModBytes, nil)
		if err != nil {
			return nil, err
		}
		abs, err := filepath.Abs(workDir)
		if err != nil {
			return nil, err
		}
		return &ModuleStruct{
			Path:      mod.Module.Mod.Path,
			Dir:       abs,
			GoVersion: mod.Go.Version,
		}, nil
	}
	// 通过 go list -json -m 获取
	ms, err := GetGoMods(workDir)
	if err != nil {
		return nil, err
	}

	if len(ms) == 0 {
		return nil, errors.New("not go module project")
	}

	// mono project
	for _, m := range ms {
		if filepath.Clean(workDir) == filepath.Clean(m.Dir) {
			return &m, nil
		}
	}

	// unknown
	return &ms[0], nil
}

func GetGoMods(workDir string) ([]ModuleStruct, error) {
	command := exec.Command("go", "list", "-json", "-m")
	command.Dir = workDir
	data, err := command.CombinedOutput()
	if err != nil {
		if strings.Contains(string(data), "go mod tidy") {
			if !config.C.Quiet {
				fmt.Printf("%s go mod tidy. Please wait...\n", console.Green("Running"))
			}
			if _, err = execx.Run("go mod tidy", workDir); err != nil {
				return nil, err
			}
			command = exec.Command("go", "list", "-json", "-m")
			command.Dir = workDir
			if data, err = command.CombinedOutput(); err != nil {
				return nil, errors.New(string(data))
			}
		} else {
			return nil, errors.New(string(data))
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
