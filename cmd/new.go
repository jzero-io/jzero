/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	new2 "github.com/jzero-io/jzero/internal/new"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/zeromicro/go-zero/core/color"

	git "github.com/go-git/go-git/v5"
	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "jzero new project",
	Long:  `jzero new project`,
	PreRun: func(_ *cobra.Command, args []string) {
		new2.Version = Version
		new2.AppName = args[0]

		if new2.Module == "" {
			new2.Module = args[0]
		}

		if new2.Output == "" {
			new2.Output = args[0]
		}

		if new2.Remote != "" && new2.Branch != "" {
			// clone to local
			home, _ := os.UserHomeDir()
			_ = os.MkdirAll(filepath.Join(home, ".jzero"), 0o755)
			if !pathx.FileExists(filepath.Join(home, ".jzero", "templates", new2.Branch)) {
				fmt.Printf("%s templates into '%s/templates/%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), filepath.Join(home, ".jzero"), new2.Branch)
				_, err := git.PlainClone(filepath.Join(home, ".jzero", "templates", new2.Branch), false, &git.CloneOptions{
					SingleBranch:  true,
					URL:           new2.Remote,
					Depth:         0,
					ReferenceName: plumbing.ReferenceName("refs/heads/" + new2.Branch),
				})
				cobra.CheckErr(err)
				fmt.Println(color.WithColor("Done", color.FgGreen))
			} else {
				fmt.Printf("%s cache: %s\n", color.WithColor("Using", color.FgGreen), filepath.Join(home, ".jzero", "templates", new2.Branch))
			}
			embeded.Home = filepath.Join(home, ".jzero", "templates", new2.Branch)
		}
	},
	RunE: new2.NewProject,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&new2.Module, "module", "m", "", "set go module")
	newCmd.Flags().StringVarP(&new2.Output, "output", "o", "", "set output dir")
	newCmd.Flags().StringVarP(&new2.AppDir, "app-dir", "", "", "set app dir")
	newCmd.Flags().StringVarP(&embeded.Home, "home", "", "", "set home dir")
	newCmd.Flags().StringVarP(&new2.ConfigType, "config-type", "", "yaml", "set config type, default toml")
	newCmd.Flags().StringVarP(&new2.Remote, "remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
	newCmd.Flags().StringVarP(&new2.Branch, "branch", "b", "", "remote templates repo branch")
}
