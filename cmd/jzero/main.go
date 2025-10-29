/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package main

import (
	"embed"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/check"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/completion"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/format"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/mcp"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/migrate"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/new"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/serverless"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/template"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/upgrade"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/version"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/hooks"
	mcppkg "github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mcp"
)

var WorkingDir string

// embeded
var (
	//go:embed all:.template
	Template embed.FS
)

// ldflags
var (
	Version = "0.64.5"
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
	Short: `Used to create project by templates and generate server/client code by api/proto/sql file.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Run environment check first
		if cmd.Name() != check.GetCommand().Use && cmd.Name() != version.GetCommand().Use {
			frameType, err := desc.GetFrameType()
			if err != nil {
				return err
			}
			if frameType != "" {
				if err := check.RunCheck(false); err != nil {
					return err
				}
			}
		}
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
	rootCmd.PersistentFlags().StringVarP(&config.CfgFile, "config", "f", ".jzero.yaml", "set config file")
	rootCmd.PersistentFlags().StringVarP(&config.CfgEnvFile, "config-env", "", ".jzero.env.yaml", "set config env file")
	rootCmd.PersistentFlags().StringSliceP("register-tpl-val", "", []string{}, "register tpl value, e.g. --register-tpl-val key=value")
	rootCmd.PersistentFlags().BoolP("debug", "", false, "debug mode")
	rootCmd.PersistentFlags().IntP("debug-sleep-time", "", 0, "debug sleep time")

	rootCmd.AddCommand(check.GetCommand())
	rootCmd.AddCommand(completion.GetCommand())
	rootCmd.AddCommand(format.GetCommand())
	rootCmd.AddCommand(gen.GetCommand())
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

	cobra.CheckErr(config.InitConfig(rootCmd))
	if config.C.Debug {
		logx.MustSetup(logx.LogConf{Encoding: "plain"})
		logx.SetLevel(logx.DebugLevel)
		if config.C.DebugSleepTime > 0 {
			logx.Debugf("using jzero frame debug mode, please wait time.Sleep(time.Second * %d)", config.C.DebugSleepTime)
		} else {
			logx.Debugf("using jzero frame debug mode")
		}
		time.Sleep(time.Duration(config.C.DebugSleepTime) * time.Second)
		logx.Debugf("get config: %#v", config.C)
	} else {
		logx.Disable()
	}
}
