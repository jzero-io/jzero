/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/internal/template/templatebuild"
	"github.com/jzero-io/jzero/internal/template/templateinit"
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "jzero template",
	Long:  `jzero template`,
}

var templateInitCmd = &cobra.Command{
	Use:   "init",
	Short: "jzero template init",
	Long:  `jzero template init`,
	PreRun: func(_ *cobra.Command, _ []string) {
		if templateinit.Home == "" {
			home, _ := os.UserHomeDir()
			templateinit.Home = filepath.Join(home, ".jzero", Version)
		}
	},
	RunE: templateinit.Init,
}

var templateBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "jzero template build",
	Long:  `jzero template build`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if templatebuild.Output == "" {
			home, _ := os.UserHomeDir()
			templatebuild.Output = filepath.Join(home, ".jzero", "templates", templatebuild.Name, "app")
		} else {
			templatebuild.Output = filepath.Join(templatebuild.Output, "app")
		}
	},
	RunE: templatebuild.Build,
}

func init() {
	rootCmd.AddCommand(templateCmd)

	{
		templateCmd.AddCommand(templateBuildCmd)

		templateBuildCmd.Flags().StringVarP(&templatebuild.WorkingDir, "working-dir", "w", ".", "default working directory")
		templateBuildCmd.Flags().StringVarP(&templatebuild.Name, "name", "n", "", "template name")
		_ = templateBuildCmd.MarkFlagRequired("name")
		templateBuildCmd.Flags().StringVarP(&templatebuild.Output, "output", "o", "", "default output directory")
	}

	{
		templateCmd.AddCommand(templateInitCmd)

		templateInitCmd.Flags().StringVarP(&templateinit.Remote, "remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
		templateInitCmd.Flags().StringVarP(&templateinit.Branch, "branch", "b", "", "remote templates repo branch")
		templateInitCmd.Flags().StringVarP(&templateinit.Home, "home", "", "", "template output dir")
	}
}
