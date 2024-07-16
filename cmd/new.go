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
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/new"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: `Used to create project from templates`,
	PreRun: func(_ *cobra.Command, args []string) {
		new.AppName = args[0]

		if new.Output == "" {
			new.Output = args[0]

			if pathx.FileExists(new.Output) {
				cobra.CheckErr(errors.Errorf("%s already exists", new.Output))
			}
		}
		if new.Module == "" {
			new.Module = filepath.ToSlash(new.Output)
		}

		if !pathx.FileExists(embeded.Home) {
			home, _ := os.UserHomeDir()
			embeded.Home = filepath.Join(home, ".jzero", Version)
		}

		if new.Remote != "" && new.Branch != "" {
			// clone to local
			home, _ := os.UserHomeDir()
			_ = os.MkdirAll(filepath.Join(home, ".jzero"), 0o755)

			if !pathx.FileExists(filepath.Join(home, ".jzero", "templates", new.Branch)) || !new.Cache {
				fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), filepath.Join(home, ".jzero", "templates", new.Branch))

				_ = os.RemoveAll(filepath.Join(home, ".jzero", "templates", new.Branch))

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				_, err := git.PlainCloneContext(ctx, filepath.Join(home, ".jzero", "templates", new.Branch), false, &git.CloneOptions{
					SingleBranch:  true,
					URL:           new.Remote,
					Depth:         0,
					ReferenceName: plumbing.ReferenceName("refs/heads/" + new.Branch),
				})
				cobra.CheckErr(err)
				fmt.Println(color.WithColor("Done", color.FgGreen))
			} else {
				fmt.Printf("%s cache: %s\n", color.WithColor("Using", color.FgGreen), filepath.Join(home, ".jzero", "templates", new.Branch))
			}
			embeded.Home = filepath.Join(home, ".jzero", "templates", new.Branch)

			if new.WithTemplate {
				fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), filepath.Join(new.Output, ".template"))
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				_, err := git.PlainCloneContext(ctx, filepath.Join(new.Output, ".template"), false, &git.CloneOptions{
					SingleBranch:  true,
					URL:           new.Remote,
					Depth:         0,
					ReferenceName: plumbing.ReferenceName("refs/heads/" + new.Branch),
				})
				cobra.CheckErr(err)
				fmt.Println(color.WithColor("Done", color.FgGreen))

				if pathx.FileExists(filepath.Join(new.Output, ".template", "go-zero")) {
					fmt.Printf("If you want to use all go-zero templates. Please exec `goctl template init --home %s/.template/go-zero`\n", new.Output)
				}
			}
		}
	},
	RunE: new.NewProject,
	Args: cobra.ExactArgs(1),
}

func init() {
	wd, _ := os.Getwd()

	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&new.Module, "module", "m", "", "set go module")
	newCmd.Flags().StringVarP(&new.Output, "output", "o", "", "set output dir")
	newCmd.Flags().StringVarP(&embeded.Home, "home", "", filepath.Join(wd, ".template"), "set home dir")
	newCmd.Flags().StringVarP(&new.Remote, "remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
	newCmd.Flags().StringVarP(&new.Branch, "branch", "b", "", "remote templates repo branch")
	newCmd.Flags().BoolVarP(&new.Cache, "cache", "", false, "get templates in local templates dir")
	newCmd.Flags().BoolVarP(&new.WithTemplate, "with-template", "", false, "with template files in your project")
	newCmd.Flags().StringVarP(&new.Style, "style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
	newCmd.Flags().StringSliceVarP(&new.Features, "features", "", []string{}, "select features")
}
