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
	PreRun: func(_ *cobra.Command, _ []string) {
		new.Version = Version
	},
	RunE: new.NewProject,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&new.Module, "module", "m", "", "set go module")
	_ = newCmd.MarkFlagRequired("module")

	newCmd.Flags().StringVarP(&new.Dir, "dir", "d", "", "set output dir")
	_ = newCmd.MarkFlagRequired("dir")

	newCmd.Flags().StringVarP(&new.APP, "app", "", "", "set app name")
	_ = newCmd.MarkFlagRequired("app")

	newCmd.Flags().StringVarP(&embeded.Home, "home", "", "", "set home dir")
}
