/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	cfgFile string
	Debug   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jzero",
	Short: "jzero framework",
	Long:  `jzero framework.`,
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jzero/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		wd, err := os.Getwd()
		cobra.CheckErr(err)

		var (
			configPath string
			configType string
			configName string
		)
		if pathx.FileExists(filepath.Join(wd, "config.toml")) {
			configPath = wd
			configType = "toml"
			configName = "config"
		} else {
			configPath = filepath.Join(home, ".jzero")
			configType = "toml"
			configName = "config"
		}

		viper.AddConfigPath(configPath)
		viper.SetConfigType(configType)
		viper.SetConfigName(configName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		cfgFile = viper.ConfigFileUsed()
	}
}
