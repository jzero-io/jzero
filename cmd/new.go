/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/jzero-io/jzero/cmd/new"
	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
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
}
