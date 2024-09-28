/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/template/templatebuild"
	"github.com/jzero-io/jzero/internal/template/templateinit"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: `Used to save and build templates`,
}

var templateInitCmd = &cobra.Command{
	Use:   "init",
	Short: `Save template files on your disk`,
	PreRun: func(_ *cobra.Command, _ []string) {
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.C.Template.Init.Output == "" {
			home, _ := os.UserHomeDir()
			config.C.Template.Init.Output = filepath.Join(home, ".jzero", Version)
		}
		return templateinit.Init(config.C)
	},
}

var templateBuildCmd = &cobra.Command{
	Use:   "build",
	Short: `Build your current project to template`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if config.C.Template.Build.Output == "" {
			home, _ := os.UserHomeDir()
			config.C.Template.Build.Output = filepath.Join(home, ".jzero", "templates", config.C.Template.Build.Name, "app")
		} else {
			config.C.Template.Build.Output = filepath.Join(config.C.Template.Build.Output, "app")
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return templatebuild.Build(config.C.Template)
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	{
		templateCmd.AddCommand(templateBuildCmd)

		templateBuildCmd.Flags().StringP("working-dir", "w", ".", "default working directory")
		templateBuildCmd.Flags().StringP("name", "n", "", "template name")
		_ = templateBuildCmd.MarkFlagRequired("name")
		templateBuildCmd.Flags().StringP("output", "o", "", "default output directory")
	}

	{
		templateCmd.AddCommand(templateInitCmd)

		templateInitCmd.Flags().StringP("remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
		templateInitCmd.Flags().StringP("branch", "b", "", "remote templates repo branch")
		templateInitCmd.Flags().StringP("output", "o", "", "template output dir")
	}
}
