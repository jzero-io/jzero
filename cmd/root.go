/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Debug   bool
	CfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "jzero",
	Short: `Used to create project by templates and generate server/client code by proto and api file.
`,
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

	rootCmd.PersistentFlags().StringVarP(&CfgFile, "config", "f", ".jzero.yaml", "set config file")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if pathx.FileExists(CfgFile) {
		viper.SetConfigFile(CfgFile)
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			cobra.CheckErr(err)
		}
	}

	if Debug {
		logx.MustSetup(logx.LogConf{Encoding: "plain"})
		logx.SetLevel(logx.DebugLevel)
		logx.Debugf("using jzero frame debug mode, please wait time.Sleep(time.Second * 10)")
		time.Sleep(time.Second * 10)
	} else {
		logx.Disable()
	}
}
