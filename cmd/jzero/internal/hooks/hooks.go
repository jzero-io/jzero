package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/execx"
)

func Run(cmd *cobra.Command, hookAction, hooksName string, hooks []string) error {
	wd, _ := os.Getwd()

	quiet, _ := strconv.ParseBool(cmd.Flags().Lookup("quiet").Value.String())

	if os.Getenv("JZERO_HOOK_TRIGGERED") == "true" {
		return nil
	}

	if os.Getenv("JZERO_FORKED") == "true" && hooksName == "global" && hookAction == "Before" {
		return nil
	}

	if len(hooks) > 0 && !quiet {
		fmt.Printf("%s\n", console.Green(fmt.Sprintf("Start %s %s hooks", hookAction, hooksName)))
	}

	for _, v := range hooks {
		if !quiet {
			fmt.Printf("%s command %s\n", console.Green("Run"), v)
		}
		err := execx.Run(v, wd, "JZERO_HOOK_TRIGGERED=true")
		if err != nil {
			return err
		}
	}

	if len(hooks) > 0 && !quiet {
		fmt.Printf("%s\n", console.Green("Done"))
	}

	// fork 一个子进程来运行后续的指令
	if len(hooks) > 0 && hookAction == "Before" && hooksName == "global" {
		logx.Debugf("Before hooks executed, forking a new process to continue")

		// 获取当前可执行文件路径
		executable, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %v", err)
		}
		// 准备命令行参数
		args := os.Args[1:]

		// 设置环境变量，防止无限递归
		env := append(os.Environ(), "JZERO_FORKED=true")

		// 创建新进程
		fork := exec.Command(executable, args...)
		fork.Env = env
		fork.Stdin = os.Stdin
		fork.Stdout = os.Stdout
		fork.Stderr = os.Stderr

		// 启动新进程
		if err = fork.Start(); err != nil {
			return fmt.Errorf("failed to start forked process: %v", err)
		}
		// 等待新进程完成
		if err = fork.Wait(); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				os.Exit(exitErr.ExitCode())
			}
			return fmt.Errorf("forked process failed: %v", err)
		}
		// 子进程成功完成，退出当前进程
		os.Exit(0)
	}

	return nil
}
