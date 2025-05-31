/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package main

import (
	"embed"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/a8m/envsubst"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"gopkg.in/yaml.v3"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/check"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/completion"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/format"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/ivm"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/mcp"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/migrate"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/new"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/serverless"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/template"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/upgrade"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/version"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/hooks"
	mcppkg "github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mcp"
)

var (
	CfgFile    string
	CfgEnvFile string
	WorkingDir string
)

// embeded
var (
	//go:embed all:.template
	Template embed.FS
)

// ldflags
var (
	Version = "0.39.1"
	Commit  string
	Date    string
)

func main() {
	embeded.Template = Template
	version.Version = Version
	version.Date = Date
	version.Commit = Commit

	os.Args = mcppkg.ProcessOsArgs()

	Execute()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "jzero",
	Short: `Used to create project by templates and generate server/client code by proto and api file.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return hooks.Run(cmd, "Before", "global", config.C.Hooks.Before)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if parseBool, err := strconv.ParseBool(cmd.Flags().Lookup("version").Value.String()); err == nil && parseBool {
			version.GetVersion()
			return
		}
		if err := cmd.Help(); err != nil {
			cobra.CheckErr(err)
		}
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return hooks.Run(cmd, "After", "global", config.C.Hooks.After)
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
	cobra.OnInitialize(InitConfig)

	rootCmd.Flags().BoolP("version", "v", false, "show version")
	rootCmd.PersistentFlags().StringVarP(&WorkingDir, "working-dir", "w", "", "set working directory")
	rootCmd.PersistentFlags().StringVarP(&CfgFile, "config", "f", ".jzero.yaml", "set config file")
	rootCmd.PersistentFlags().StringVarP(&CfgEnvFile, "config-env", "", ".jzero.env.yaml", "set config env file")
	rootCmd.PersistentFlags().BoolP("debug", "", false, "debug mode")
	rootCmd.PersistentFlags().IntP("debug-sleep-time", "", 0, "debug sleep time")

	rootCmd.AddCommand(check.GetCommand())
	rootCmd.AddCommand(completion.GetCommand())
	rootCmd.AddCommand(format.GetCommand())
	rootCmd.AddCommand(gen.GetCommand())
	rootCmd.AddCommand(ivm.GetCommand())
	rootCmd.AddCommand(mcp.GetCommand())
	rootCmd.AddCommand(migrate.GetCommand())
	rootCmd.AddCommand(new.GetCommand())
	rootCmd.AddCommand(serverless.GetCommand())
	rootCmd.AddCommand(template.GetCommand())
	rootCmd.AddCommand(upgrade.GetCommand())
	rootCmd.AddCommand(version.GetCommand())
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if len(os.Args) >= 2 {
		if os.Args[1] == version.GetCommand().Use {
			return
		}
	}

	if WorkingDir != "" {
		if err := os.Chdir(WorkingDir); err != nil {
			cobra.CheckErr(err)
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

	if err := config.TraverseCommands("", rootCmd); err != nil {
		panic(err)
	}

	if config.C.Debug {
		logx.MustSetup(logx.LogConf{Encoding: "plain"})
		logx.SetLevel(logx.DebugLevel)
		if config.C.DebugSleepTime > 0 {
			logx.Debugf("using jzero frame debug mode, please wait time.Sleep(time.Second * %d)", config.C.DebugSleepTime)
		} else {
			logx.Debugf("using jzero frame debug mode")
		}
		time.Sleep(time.Duration(config.C.DebugSleepTime) * time.Second)
	} else {
		logx.Disable()
	}
}
