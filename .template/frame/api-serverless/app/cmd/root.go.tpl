package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CfgFile string

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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "etc/etc.yaml", "config file (default is project root dir etc/etc.yaml")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(os.Args) <= 1 || os.Args[1] != serverCmd.Name() {
		return
	}

    viper.SetConfigFile(CfgFile)
	viper.AutomaticEnv() // read in environment variables that match
}