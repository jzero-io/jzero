/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>
*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/serverless/serverlessbuild"
	"github.com/jzero-io/jzero/internal/serverless/serverlessdelete"
)

// serverlessCmd represents the serverless command
var serverlessCmd = &cobra.Command{
	Use:   "serverless",
	Short: "jzero serverless",
	Long:  `jzero serverless.`,
}

var serverlessBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "jzero serverless build",
	Long:  `jzero serverless build.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		embeded.Home = config.C.Serverless.Home
		return serverlessbuild.Run()
	},
}

var serverlessDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "jzero serverless delete",
	Long:  `jzero serverless delete.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		embeded.Home = config.C.Serverless.Home
		return serverlessdelete.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverlessCmd)
	serverlessCmd.AddCommand(serverlessBuildCmd)
	serverlessCmd.AddCommand(serverlessDeleteCmd)

	serverlessCmd.PersistentFlags().StringP("home", "", ".template", "set templates path")
	serverlessDeleteCmd.Flags().StringSliceP("plugin", "p", nil, "plugin name")
}
