package cmd

import (
	"github.com/spf13/cobra"

	"{{ .Module }}/daemon"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "{{ .APP }} daemon",
	Long:  "{{ .APP }} daemon",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Start(cfgFile)
		select {}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
