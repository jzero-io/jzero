/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jaronnie/jzero/cmd/new"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "jzero new project",
	Long:  `jzero new project`,
	RunE:  new.NewProject,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&new.Module, "module", "m", "", "set go module")
	newCmd.Flags().StringVarP(&new.Dir, "dir", "d", "", "set output dir")
	newCmd.Flags().StringVarP(&new.APP, "app", "", "", "set app name")
}
