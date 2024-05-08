package mod

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

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
		return nil, nil
	}

	var m ModuleStruct
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
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
