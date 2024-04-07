/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jaronnie/worktab/worktabd"
)

// worktabdCmd represents the worktabd command
var worktabdCmd = &cobra.Command{
	Use:   "worktabd",
	Short: "worktabd daemon",
	Long:  `worktabd daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		worktabd.StartWorktabDaemon(cfgFile)

		select {}
	},
}

func init() {
	rootCmd.AddCommand(worktabdCmd)
}
