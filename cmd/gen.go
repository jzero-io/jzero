/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/jzero-io/jzero/cmd/gen"
	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "jzero gen code",
	Long:  `jzero gen code`,
	PreRun: func(_ *cobra.Command, _ []string) {
		gen.Version = Version
	},
	RunE: gen.Gen,
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVarP(&gen.WorkingDir, "working-dir", "w", "", "set working dir")

	dir, _ := os.UserHomeDir()
	genCmd.Flags().StringVarP(&embeded.Home, "home", "", filepath.Join(dir, ".jzero"), "set template home")
}
