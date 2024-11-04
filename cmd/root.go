/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/a8m/envsubst"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"gopkg.in/yaml.v3"

	"github.com/jzero-io/jzero/config"
)

var (
	CfgFile    string
	CfgEnvFile string
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
	rootCmd.PersistentFlags().StringVarP(&CfgEnvFile, "config-env", "", ".jzero.env.yaml", "set config env file")
	rootCmd.PersistentFlags().BoolP("debug", "", false, "debug mode")
	rootCmd.PersistentFlags().IntP("debug-sleep-time", "", 3, "debug sleep time")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(os.Args) >= 2 {
		if os.Args[1] == versionCmd.Use {
			return
		}
	}

	if pathx.FileExists(CfgFile) {
		viper.SetConfigFile(CfgFile)
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			cobra.CheckErr(err)
		}
	}

	if pathx.FileExists(CfgEnvFile) {
		data, err := envsubst.ReadFile(CfgEnvFile)
		if err != nil {
			log.Fatalf("envsubst error: %v", err)
		}
		var env map[string]any
		err = yaml.Unmarshal(data, &env)
		if err != nil {
			log.Fatalf("yaml unmarshal error: %v", err)
		}

		for k, v := range env {
			_ = os.Setenv(k, cast.ToString(v))
		}
	}

	if err := traverseCommands("", rootCmd); err != nil {
		panic(err)
	}

	// patch jzero version
	if config.C.Version != "" {
		fmt.Printf("use jzero version: %s\n", config.C.Version)
		if err := golang.Install(fmt.Sprintf("github.com/jzero-io/jzero@%s", config.C.Version)); err != nil {
			cobra.CheckErr(err)
		}
	}

	if config.C.Debug {
		logx.MustSetup(logx.LogConf{Encoding: "plain"})
		logx.SetLevel(logx.DebugLevel)
		logx.Debugf("using jzero frame debug mode, please wait time.Sleep(time.Second * %d)", config.C.DebugSleepTime)
		time.Sleep(time.Duration(config.C.DebugSleepTime) * time.Second)
	} else {
		logx.Disable()
	}
}

func traverseCommands(prefix string, cmd *cobra.Command) error {
	err := config.SetConfig(prefix, cmd.Flags())
	if err != nil {
		return err
	}

	for _, subCommand := range cmd.Commands() {
		newPrefix := fmt.Sprintf("%s.%s", prefix, subCommand.Use)
		if prefix == "" {
			newPrefix = subCommand.Use
		}
		err = traverseCommands(newPrefix, subCommand)
		if err != nil {
			return err
		}
	}

	return nil
}
