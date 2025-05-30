package execx

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

func Run(arg, dir string, env ...string) error {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case vars.OsMac, vars.OsLinux:
		cmd = exec.Command("sh", "-c", arg)
	case vars.OsWindows:
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return fmt.Errorf("unexpected os: %v", goos)
	}

	if len(dir) > 0 {
		cmd.Dir = dir
	}

	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		_, _ = io.Copy(os.Stdout, stdout)
	}()

	go func() {
		_, _ = io.Copy(os.Stderr, stderr)
	}()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
