/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>
*/

package serverless

import (
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/serverless/serverlessbuild"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/serverless/serverlessdelete"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
)

// serverlessCmd represents the serverless command
var serverlessCmd = &cobra.Command{
	Use:   "serverless",
	Short: "build serverless functions",
}

var serverlessBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "jzero serverless build",
	Long:  `jzero serverless build.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		embeded.Home = config.C.Home
		return serverlessbuild.Run()
	},
}

var serverlessDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "jzero serverless delete",
	Long:  `jzero serverless delete.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		embeded.Home = config.C.Home
		return serverlessdelete.Run()
	},
}

func GetCommand() *cobra.Command {
	serverlessCmd.AddCommand(serverlessBuildCmd)
	serverlessCmd.AddCommand(serverlessDeleteCmd)

	serverlessDeleteCmd.Flags().StringSliceP("plugin", "p", nil, "plugin name")

	return serverlessCmd
}
