/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/cmd/new"
	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "jzero new project",
	Long:  `jzero new project`,
	PreRun: func(_ *cobra.Command, args []string) {
		new.Version = Version
		new.APP = args[0]

		if new.Module == "" {
			new.Module = args[0]
		}

		if new.Dir == "" {
			new.Dir = args[0]
		}

		if new.Remote != "" && new.Branch != "" {
			// clone to local
			home, _ := os.UserHomeDir()
			_ = os.MkdirAll(filepath.Join(home, ".jzero"), 0o755)

			if !pathx.FileExists(filepath.Join(home, ".jzero", "templates", new.Branch)) {
				fmt.Printf("Cloning into '%s/templates/%s', please wait...\n", filepath.Join(home, ".jzero"), new.Branch)
				_, err := execx.Run(fmt.Sprintf("git clone %s -b %s templates/%s", new.Remote, new.Branch, new.Branch), filepath.Join(home, ".jzero"))
				cobra.CheckErr(err)
				fmt.Println("Clone success")
			} else {
				fmt.Printf("Using cache: %s\n", filepath.Join(home, ".jzero", "templates", new.Branch))
			}

			embeded.Home = filepath.Join(home, ".jzero", "templates", new.Branch)
		}
	},
	RunE: new.NewProject,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&new.Module, "module", "m", "", "set go module")
	newCmd.Flags().StringVarP(&new.Dir, "dir", "d", "", "set output dir")
	newCmd.Flags().StringVarP(&embeded.Home, "home", "", "", "set home dir")
	newCmd.Flags().StringVarP(&new.ConfigType, "config-type", "", "yaml", "set config type, default toml")
	newCmd.Flags().StringVarP(&new.Remote, "remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
	newCmd.Flags().StringVarP(&new.Branch, "branch", "b", "", "remote templates repo branch")
}
