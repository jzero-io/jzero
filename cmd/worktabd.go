/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jaronnie/jzero/jzerod"
)

// worktabdCmd represents the jzerod command
var worktabdCmd = &cobra.Command{
	Use:   "jzerod",
	Short: "jzerod daemon",
	Long:  `jzerod daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		jzerod.StartWorktabDaemon(cfgFile)

		select {}
	},
}

func init() {
	rootCmd.AddCommand(worktabdCmd)
}
