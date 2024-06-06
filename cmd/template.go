/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/internal/template/templatebuild"

	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
)

var Home string

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
	Run: func(_ *cobra.Command, _ []string) {
		dir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		if Home == "" {
			Home = filepath.Join(dir, ".jzero", Version)
		}

		err = embeded.WriteTemplateDir(filepath.Join(""), Home)
		cobra.CheckErr(err)
	},
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

	templateCmd.AddCommand(templateBuildCmd)

	templateBuildCmd.Flags().StringVarP(&templatebuild.WorkingDir, "working-dir", "w", ".", "default working directory")
	templateBuildCmd.Flags().StringVarP(&templatebuild.Name, "name", "n", "", "template name")
	_ = templateBuildCmd.MarkFlagRequired("name")

	templateBuildCmd.Flags().StringVarP(&templatebuild.Output, "output", "o", "", "default output directory")

	templateCmd.AddCommand(templateInitCmd)
	templateInitCmd.Flags().StringVarP(&Home, "home", "", "", "template home directory")
}
