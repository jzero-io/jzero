/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jzero-io/jzero-contrib/filex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/internal/new"
	"github.com/jzero-io/jzero/pkg/mod"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:          "new",
	Short:        `Used to create project from templates`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.C.New.Output == "" {
			config.C.New.Output = args[0]

			if pathx.FileExists(config.C.New.Output) {
				return errors.Errorf("%s already exists", config.C.New.Output)
			}
		}

		fmt.Printf("%s project %s in %s dir\n", color.WithColor("Creating", color.FgGreen), config.C.New.Output, config.C.New.Output)

		if config.C.New.Module == "" {
			config.C.New.Module = args[0]
		}
		// 在 go.mod 项目下但是项目本身没有 go.mod 文件
		if config.C.New.Mono {
			wd, _ := os.Getwd()
			parentPackage, err := mod.GetParentPackage(wd)
			if err != nil {
				return err
			}
			config.C.New.Module = filepath.ToSlash(filepath.Join(parentPackage, config.C.New.Output))
		}

		home, _ := os.UserHomeDir()

		var base string
		switch {
		// 指定特定路径作为模板
		case config.C.New.Home != "":
			embeded.Home = config.C.New.Home
			if config.C.New.Frame != "" {
				base = filepath.Join("frame", config.C.New.Frame, "app")
			} else {
				base = filepath.Join("app")
				if !pathx.FileExists(base) && pathx.FileExists(filepath.Join(embeded.Home, "frame")) {
					base = filepath.Join("frame", "api", "app")
				}
			}
		// 指定本地路径 ~/.jzero/templates/local 下的某文件夹作为模板
		case config.C.New.Local != "":
			embeded.Home = filepath.Join(home, ".jzero", "templates", "local", config.C.New.Local)
			base = filepath.Join("app")
		// 使用内置模板
		case config.C.New.Frame != "":
			base = filepath.Join("frame", config.C.New.Frame, "app")
		// 使用远程仓库模板
		case config.C.New.Remote != "" && config.C.New.Branch != "":
			fp := filepath.Join(home, ".jzero", "templates", "remote", config.C.New.Branch)
			if filex.DirExists(fp) && config.C.New.Cache {
				fmt.Printf("%s cache templates from '%s', please wait...\n", color.WithColor("Using", color.FgGreen), fp)
			} else {
				_ = os.RemoveAll(fp)
				fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), fp)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				if strings.HasPrefix(config.C.New.Remote, "git@") {
					// SSH 协议
					commandContext := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--branch", config.C.New.Branch, config.C.New.Remote, fp)
					if resp, err := commandContext.CombinedOutput(); err != nil {
						return errors.New(string(resp))
					}
				} else {
					// HTTP 协议
					auth := &http.BasicAuth{
						Username: config.C.New.RemoteAuthUsername,
						Password: config.C.New.RemoteAuthPassword,
					}
					// clone to local
					if _, err := git.PlainCloneContext(ctx, fp, false, &git.CloneOptions{
						SingleBranch:  true,
						URL:           config.C.New.Remote,
						Depth:         0,
						ReferenceName: plumbing.ReferenceName("refs/heads/" + config.C.New.Branch),
						Auth:          auth,
					}); err != nil {
						return err
					}
				}

				_ = os.RemoveAll(filepath.Join(fp, ".git"))
			}
			fmt.Println(color.WithColor("Done", color.FgGreen))
			embeded.Home = fp
			base = filepath.Join("app")
		default:
			// 默认使用 api 模板
			config.C.New.Frame = "api"
			base = filepath.Join("frame", "api", "app")
		}

		// 没有设置 home，则使用内置模板持久化的默认路径 ~/.jzero/templates/$version
		if !pathx.FileExists(embeded.Home) {
			embeded.Home = filepath.Join(home, ".jzero", "templates", Version)
		}
		return new.Run(args[0], base)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		if err = os.Chdir(config.C.New.Output); err != nil {
			return err
		}
		defer func() {
			if err = os.Chdir(dir); err != nil {
				cobra.CheckErr(err)
			}
		}()

		// special dir for jzero
		if !filex.DirExists("desc") {
			return nil
		}
		fmt.Printf("%s desc dir in %s, auto generate code\n", color.WithColor("Detected", color.FgGreen), config.C.New.Output)

		initConfig()
		// for gen persistent flags
		if config.C.Gen.Style == "" {
			config.C.Gen.Style = "gozero"
		}
		if config.C.Gen.Home == "" {
			config.C.Gen.Home = filepath.Join(config.C.Wd(), ".template")
		}

		return gen.Run(true)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("module", "m", "", "set go module")
	newCmd.Flags().StringP("output", "o", "", "set output dir")
	newCmd.Flags().StringP("home", "", "", "use the specified template.")
	newCmd.Flags().StringP("frame", "", "", "set frame")
	newCmd.Flags().StringP("remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
	newCmd.Flags().StringP("remote-auth-username", "", "", "remote templates repo auth username")
	newCmd.Flags().StringP("remote-auth-password", "", "", "remote templates repo auth password")
	newCmd.Flags().StringP("branch", "b", "", "use remote template repo branch")
	newCmd.Flags().BoolP("cache", "", false, "remote template using cache")
	newCmd.Flags().StringP("local", "", "", "use local template")
	newCmd.Flags().StringSliceP("features", "", []string{}, "select features")
	newCmd.Flags().BoolP("mono", "", false, "mono project under go mod project")
}
