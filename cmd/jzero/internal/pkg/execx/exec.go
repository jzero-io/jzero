package execx

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

// RunOutput runs a command and returns its output
func RunOutput(arg, dir string, env ...string) (string, error) {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case vars.OsMac, vars.OsLinux:
		cmd = exec.Command("sh", "-c", arg)
	case vars.OsWindows:
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return "", fmt.Errorf("unexpected os: %v", goos)
	}

	if len(dir) > 0 {
		cmd.Dir = dir
	}

	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}
