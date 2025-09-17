/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package template

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/template/templatebuild"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/template/templateinit"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/version"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
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
				config.C.Template.Init.Output = filepath.Join(home, ".jzero", "templates", version.Version)
			}
		}
		return templateinit.Run()
	},
}

var templateBuildCmd = &cobra.Command{
	Use:   "build",
	Short: `Build your current project to template`,
	Long:  `Build your current project to template and save them into ${HOME}/.jzero/templates/local.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.C.Template.Build.Output == "" {
			home, _ := os.UserHomeDir()
			config.C.Template.Build.Output = filepath.Join(home, ".jzero", "templates", "local", config.C.Template.Build.Name)
		}
		return templatebuild.Run(config.C.Template)
	},
}

func GetCommand() *cobra.Command {
	{
		templateCmd.AddCommand(templateBuildCmd)

		templateBuildCmd.Flags().StringP("working-dir", "w", ".", "working directory")
		templateBuildCmd.Flags().StringP("name", "n", "", "template name")
		_ = templateBuildCmd.MarkFlagRequired("name")
		templateBuildCmd.Flags().StringP("output", "o", "", "output directory")
		templateBuildCmd.Flags().StringSliceP("ignore", "i", []string{".git", ".idea", ".vscode", ".DS_Store", "node_modules"}, "dir list for ignored files")
	}

	{
		templateCmd.AddCommand(templateInitCmd)

		templateInitCmd.Flags().StringP("remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
		templateInitCmd.Flags().StringP("branch", "b", "", "remote template repo branch. If not set, init the embedded templates.")
		templateInitCmd.Flags().StringP("output", "o", "", "template output dir")
	}

	return templateCmd
}
