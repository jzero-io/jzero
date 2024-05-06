/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"os"
	"path/filepath"

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

		err = embeded.WriteTemplateDir(filepath.Join("go-zero"), filepath.Join(Home, "go-zero"))
		cobra.CheckErr(err)

		err = embeded.WriteTemplateDir(filepath.Join("jzero"), filepath.Join(Home, "jzero"))
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateInitCmd)

	templateInitCmd.Flags().StringVarP(&Home, "home", "", "", "template home directory")
}
