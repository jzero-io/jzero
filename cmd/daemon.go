/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jaronnie/jzero/daemon"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "jzero daemon",
	Long:  `jzero daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Start(cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
