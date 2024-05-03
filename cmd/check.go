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
		// install goctl
		_, err := LookUpTool("goctl")
		if err != nil {
			err = golang.Install("github.com/zeromicro/go-zero/tools/goctl@latest")
			cobra.CheckErr(err)
		}

		// goctl env check
		resp, err := execx.Run("goctl env check --install --verbose --force", "")
		cobra.CheckErr(err)
		fmt.Println(resp)

		log := console.NewColorConsole(true)

		// jzero env check
		_, err = LookUpTool("goreleaser")
		if err != nil {
			log.Info(`[jzero-env]: looking up "goreleaser"`)
			err = golang.Install("github.com/goreleaser/goreleaser@latest")
			cobra.CheckErr(err)
		}
		if _, err = LookUpTool("goreleaser"); err == nil {
			fmt.Println()
			log.Success(`[jzero-env]: "goreleaser" is installed`)
		}

		_, err = LookUpTool("task")
		if err != nil {
			err = golang.Install("github.com/go-task/task/v3/cmd/task@latest")
			cobra.CheckErr(err)
		}
		if _, err = LookUpTool("task"); err == nil {
			log.Success(`[jzero-env]: "task" is installed`)
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
