package cmd

import (
	"github.com/spf13/cobra"

	"{{ .Module }}/app"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "{{ .APP }} server",
	Long:  "{{ .APP }} server",
	Run: func(cmd *cobra.Command, args []string) {
		app.Start(cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
