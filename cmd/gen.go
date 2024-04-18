/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/jaronnie/jzero/cmd/gen"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "jzero gen code",
	Long:  `jzero gen code`,
	RunE:  gen.Gen,
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVarP(&gen.WorkingDir, "working-dir", "w", "", "set working dir")
}
