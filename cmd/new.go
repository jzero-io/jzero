/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/new"
	"github.com/jzero-io/jzero/pkg/mod"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: `Used to create project from templates`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.C.New.Output == "" {
			config.C.New.Output = args[0]

			if pathx.FileExists(config.C.New.Output) {
				cobra.CheckErr(errors.Errorf("%s already exists", config.C.New.Output))
			}
		}
		if config.C.New.Module == "" {
			config.C.New.Module = args[0]
		}
		// 在 go.mod 项目下但是项目本身没有 go.mod 文件
		if config.C.New.Mono {
			wd, _ := os.Getwd()
			var err error
			parentPackage, err := mod.GetParentPackage(wd)
			config.C.New.Module = filepath.ToSlash(filepath.Join(parentPackage, config.C.New.Output))
			cobra.CheckErr(err)
		}

		home, _ := os.UserHomeDir()

		// 使用远程仓库模板
		if config.C.New.Remote != "" && config.C.New.Branch != "" {
			// clone to local
			fp := filepath.Join(home, ".jzero", "templates", "remote", config.C.New.Branch)
			_ = os.MkdirAll(fp, 0o755)
			fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), fp)
			_ = os.RemoveAll(fp)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			_, err := git.PlainCloneContext(ctx, fp, false, &git.CloneOptions{
				SingleBranch:  true,
				URL:           config.C.New.Remote,
				Depth:         0,
				ReferenceName: plumbing.ReferenceName("refs/heads/" + config.C.New.Branch),
				Auth: &http.BasicAuth{
					Username: os.Getenv("JZERO_REMOTE_USERNAME"),
					Password: os.Getenv("JZERO_REMOTE_PASSWORD"),
				},
			})
			cobra.CheckErr(err)
			_ = os.RemoveAll(filepath.Join(fp, ".git"))
			fmt.Println(color.WithColor("Done", color.FgGreen))
			embeded.Home = fp
		}

		// 使用本地模板
		if config.C.New.Local != "" {
			embeded.Home = filepath.Join(home, ".jzero", "templates", "local", config.C.New.Local)
		}

		// 指定 home 时优先级最高
		if config.C.New.Home != "" {
			embeded.Home = config.C.New.Home
		}

		if !pathx.FileExists(embeded.Home) {
			embeded.Home = filepath.Join(home, ".jzero", "templates", Version)
		}
		return new.Run(config.C, args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("module", "m", "", "set go module")
	newCmd.Flags().StringP("output", "o", "", "set output dir")
	newCmd.Flags().StringP("home", "", "", "use the specified template.")
	newCmd.Flags().StringP("frame", "", "api", "frame")
	newCmd.Flags().StringP("remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
	newCmd.Flags().StringP("branch", "b", "", "use remote template repo branch")
	newCmd.Flags().StringP("local", "", "", "use local template")
	newCmd.Flags().StringP("style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
	newCmd.Flags().StringSliceP("features", "", []string{}, "select features")
	newCmd.Flags().BoolP("mono", "", false, "mono project under go mod project")
}
