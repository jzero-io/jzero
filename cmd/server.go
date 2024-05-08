/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/app"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "jzero server",
	Long:  `jzero server`,
	Run: func(_ *cobra.Command, _ []string) {
		app.Start(cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
