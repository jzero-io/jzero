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
		// 在 go mod 项目的一个子模块
		if config.C.New.SubModule {
			wd, _ := os.Getwd()
			var err error
			parentPackage, err := mod.GetParentPackage(wd)
			config.C.New.Module = filepath.ToSlash(filepath.Join(parentPackage, config.C.New.Output))
			cobra.CheckErr(err)
		}

		if !pathx.FileExists(config.C.New.Home) {
			home, _ := os.UserHomeDir()
			config.C.New.Home = filepath.Join(home, ".jzero", Version)
		}

		if config.C.New.Remote != "" && config.C.New.Branch != "" {
			// clone to local
			home, _ := os.UserHomeDir()
			_ = os.MkdirAll(filepath.Join(home, ".jzero"), 0o755)

			if !pathx.FileExists(filepath.Join(home, ".jzero", "templates", config.C.New.Branch)) || !config.C.New.Cache {
				fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), filepath.Join(home, ".jzero", "templates", config.C.New.Branch))

				_ = os.RemoveAll(filepath.Join(home, ".jzero", "templates", config.C.New.Branch))

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				_, err := git.PlainCloneContext(ctx, filepath.Join(home, ".jzero", "templates", config.C.New.Branch), false, &git.CloneOptions{
					SingleBranch:  true,
					URL:           config.C.New.Remote,
					Depth:         0,
					ReferenceName: plumbing.ReferenceName("refs/heads/" + config.C.New.Branch),
				})
				cobra.CheckErr(err)
				fmt.Println(color.WithColor("Done", color.FgGreen))
			} else {
				fmt.Printf("%s cache: %s\n", color.WithColor("Using", color.FgGreen), filepath.Join(home, ".jzero", "templates", config.C.New.Branch))
			}
			config.C.New.Home = filepath.Join(home, ".jzero", "templates", config.C.New.Branch)

			if config.C.New.WithTemplate {
				fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), filepath.Join(config.C.New.Output, ".template"))
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				_, err := git.PlainCloneContext(ctx, filepath.Join(config.C.New.Output, ".template"), false, &git.CloneOptions{
					SingleBranch:  true,
					URL:           config.C.New.Remote,
					Depth:         0,
					ReferenceName: plumbing.ReferenceName("refs/heads/" + config.C.New.Branch),
				})
				cobra.CheckErr(err)
				fmt.Println(color.WithColor("Done", color.FgGreen))

				if pathx.FileExists(filepath.Join(config.C.New.Output, ".template", "go-zero")) {
					fmt.Printf("If you want to use all go-zero templates. Please exec `goctl template init --home %s/.template/go-zero`\n", config.C.New.Output)
				}
			}
		}

		if !pathx.FileExists(config.C.New.Home) {
			home, _ := os.UserHomeDir()
			config.C.New.Home = filepath.Join(home, ".jzero", Version)
		}
		embeded.Home = config.C.New.Home

		return new.NewProject(config.C, args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	wd, _ := os.Getwd()

	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("module", "m", "", "set go module")
	newCmd.Flags().StringP("output", "o", "", "set output dir")
	newCmd.Flags().StringP("home", "", filepath.Join(wd, ".template"), "set home dir")
	newCmd.Flags().StringP("remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
	newCmd.Flags().StringP("branch", "b", "", "remote templates repo branch")
	newCmd.Flags().BoolP("cache", "", false, "get templates in local templates dir")
	newCmd.Flags().BoolP("with-template", "", false, "with template files in your project")
	newCmd.Flags().StringP("style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
	newCmd.Flags().StringSliceP("features", "", []string{}, "select features")
	newCmd.Flags().BoolP("submodule", "", false, "is project's submodule")
}
