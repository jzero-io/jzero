package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
    "github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .APP }}",
	Short: "{{ .APP }} root",
	Long:  "{{ .APP }} root.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is project root dir config.{{ .ConfigType }} or $HOME/.{{ .APP }}/config.{{ .ConfigType }})")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(os.Args) <= 1 || os.Args[1] != serverCmd.Name() {
		return
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

        wd, err := os.Getwd()
        cobra.CheckErr(err)

        var configPath string
        if pathx.FileExists(filepath.Join(wd, "config.{{ .ConfigType }}")) {
            configPath = wd
        } else {
            configPath = filepath.Join(home, ".{{ .APP }}")
        }

        viper.AddConfigPath(configPath)
        viper.SetConfigType("{{ .ConfigType }}")
        viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		cfgFile = viper.ConfigFileUsed()
	} else {
		cobra.CheckErr(err)
	}
}