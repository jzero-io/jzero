/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "jzero env check",
	Long:  `jzero env check.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := console.NewColorConsole(true)

		log.Info("[jzero-env]: looking up goctl")
		// install goctl
		_, err := LookUpTool("goctl")
		if err != nil {
			log.Warning(`[jzero-env]: goctl is not found in PATH`)
			err = golang.Install("github.com/zeromicro/go-zero/tools/goctl@latest")
			cobra.CheckErr(err)
		}
		if _, err = LookUpTool("goctl"); err == nil {
			log.Success(`[jzero-env]: "goctl" is installed`)
		} else {
			log.Fatalln("[jzero-env]: env check failed, goctl is not installed")
		}

		// goctl env check
		resp, err := execx.Run("goctl env check --install --verbose --force", "")
		cobra.CheckErr(err)
		fmt.Println(resp)

		// jzero env check
		log.Info("\n[jzero-env]: looking up task")
		_, err = LookUpTool("task")
		if err != nil {
			_ = golang.Install("github.com/go-task/task/v3/cmd/task@latest")
		}
		if _, err = LookUpTool("task"); err == nil {
			log.Success(`[jzero-env]: "task" is installed`)
		} else {
			log.Warning("[jzero-env] warning: env check failed, task is not installed")
		}

		log.Success("[jzero-env]: congratulations! your jzero environment is ready!")
	},
}

func LookUpTool(name string) (string, error) {
	suffix := getExeSuffix()
	xProtoc := name + suffix
	return env.LookPath(xProtoc)
}

func getExeSuffix() string {
	if runtime.GOOS == vars.OsWindows {
		return ".exe"
	}
	return ""
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
