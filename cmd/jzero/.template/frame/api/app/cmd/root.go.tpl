package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .APP }}",
	Short: "{{ .APP }} root",
	Long:  "{{ .APP }} root.",
    CompletionOptions: cobra.CompletionOptions{
        DisableDefaultCmd: true,
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "etc/etc.yaml", "config file (default is project root dir etc/etc.yaml")
}