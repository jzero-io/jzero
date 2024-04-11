/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jaronnie/jzero/jzerod"
)

// jzerodCmd represents the jzerod command
var jzerodCmd = &cobra.Command{
	Use:   "jzerod",
	Short: "jzero daemon",
	Long:  `jzero daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		jzerod.StartJzeroDaemon(cfgFile)

		select {}
	},
}

func init() {
	rootCmd.AddCommand(jzerodCmd)
}
