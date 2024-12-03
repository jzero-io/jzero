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
	Short: `Used to initialize or build templates`,
}

var templateInitCmd = &cobra.Command{
	Use:    "init",
	Short:  `Initialize templates`,
	Long:   `Initialize specific remote template or embedded templates on your disk`,
	PreRun: func(_ *cobra.Command, _ []string) {},

	RunE: func(cmd *cobra.Command, args []string) error {
		if config.C.Template.Init.Output == "" {
			home, _ := os.UserHomeDir()
			if config.C.Template.Init.Remote != "" && config.C.Template.Init.Branch != "" {
				config.C.Template.Init.Output = filepath.Join(home, ".jzero", "templates", "remote")
			} else {
				config.C.Template.Init.Output = filepath.Join(home, ".jzero", "templates", Version)
			}
		}
		return templateinit.Run(config.C)
	},
}

var templateBuildCmd = &cobra.Command{
	Use:   "build",
	Short: `Build your current project to template`,
	Long:  `Build your current project to template and save them in to ${HOME}/.jzero/templates/local.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if config.C.Template.Build.Output == "" {
			home, _ := os.UserHomeDir()
			config.C.Template.Build.Output = filepath.Join(home, ".jzero", "templates", "local", config.C.Template.Build.Name)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return templatebuild.Run(config.C.Template)
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	{
		templateCmd.AddCommand(templateBuildCmd)

		templateBuildCmd.Flags().StringP("working-dir", "w", ".", "working directory")
		templateBuildCmd.Flags().StringP("name", "n", "", "template name")
		_ = templateBuildCmd.MarkFlagRequired("name")
		templateBuildCmd.Flags().StringP("output", "o", "", "output directory")
		templateBuildCmd.Flags().StringSliceP("ignore", "i", templatebuild.IgnoreDirs, "dir list for ignored files")
	}

	{
		templateCmd.AddCommand(templateInitCmd)

		templateInitCmd.Flags().StringP("remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
		templateInitCmd.Flags().StringP("branch", "b", "", "remote template repo branch. If not set, init the embedded templates.")
		templateInitCmd.Flags().StringP("output", "o", "", "template output dir")
	}
}
