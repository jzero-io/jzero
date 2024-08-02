/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/spf13/cobra"
)

var Debug bool

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

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if Debug {
		logx.MustSetup(logx.LogConf{Encoding: "plain"})
		logx.SetLevel(logx.DebugLevel)
		logx.Debugf("using jzero frame debug mode, please wait time.Sleep(time.Second * 10)")
		time.Sleep(time.Second * 10)
	} else {
		logx.Disable()
	}
}
