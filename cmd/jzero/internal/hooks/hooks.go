package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

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
		var icon string
		var title string

		if hookAction == "Before" && hooksName == "global" {
			icon = ""
			title = console.Green("Executing") + " " + console.Yellow("Before Global Hooks")
		} else if hookAction == "After" && hooksName == "global" {
			icon = ""
			title = console.Green("Executing") + " " + console.Yellow("After Global Hooks")
		} else if hookAction == "Before" {
			icon = ""
			capitalName := strings.ToUpper(hooksName[:1]) + hooksName[1:]
			title = console.Green("Executing") + " " + console.Yellow("Before "+capitalName+" Command Hooks")
		} else if hookAction == "After" {
			icon = ""
			capitalName := strings.ToUpper(hooksName[:1]) + hooksName[1:]
			title = console.Green("Executing") + " " + console.Yellow("After "+capitalName+" Command Hooks")
		}

		fmt.Printf("%s\n", console.BoxHeader(icon, title))
	}

	for _, v := range hooks {
		output, err := execx.RunOutput(v, wd, "JZERO_HOOK_TRIGGERED=true")
		if !quiet {
			printHookCommand(v, err == nil)
		}
		if err != nil {
			lines := console.NormalizeErrorLines(output)
			if len(lines) == 0 {
				lines = console.NormalizeErrorLines(err.Error())
			}
			if !quiet {
				for _, line := range lines {
					fmt.Printf("%s\n", console.BoxDetailItem(line))
				}
			}
			if !quiet {
				fmt.Printf("%s\n\n", console.BoxErrorFooter())
			}
			if quiet {
				return err
			}
			return console.MarkRenderedError(err)
		}
		if !quiet && output != "" {
			printHookOutput(output)
		}
	}

	if len(hooks) > 0 && !quiet {
		fmt.Printf("%s\n\n", console.BoxSuccessFooter())
	}

	// fork 一个子进程来运行后续的指令
	if len(hooks) > 0 && hookAction == "Before" && hooksName == "global" {
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

func printHookCommand(command string, success bool) {
	if strings.Contains(command, "\n") {
		lines := strings.Split(command, "\n")
		if success {
			fmt.Printf("%s\n", console.BoxItem(console.Cyan("Executing")))
		} else {
			fmt.Printf("%s\n", console.BoxErrorItem(console.Cyan("Executing")))
		}
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine != "" {
				fmt.Printf("│  │  %s\n", trimmedLine)
			}
		}
		return
	}

	item := fmt.Sprintf("%s %s", console.Cyan("Executing"), command)
	if success {
		fmt.Printf("%s\n", console.BoxItem(item))
		return
	}

	fmt.Printf("%s\n", console.BoxErrorItem(item))
}

func printHookOutput(output string) {
	lines := strings.Split(strings.TrimRight(output, "\r\n"), "\n")
	if len(lines) == 0 {
		return
	}

	fmt.Printf("│  ╭─ %s\n", console.Cyan("Output"))
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			fmt.Print("│  │\n")
			continue
		}
		fmt.Printf("│  │  %s\n", line)
	}
	fmt.Printf("│  ╰─ %s\n", console.Cyan("Complete"))
}
