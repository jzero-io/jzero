/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>
*/

package serverless

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/serverless/serverlessbuild"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/serverless/serverlessdelete"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
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
		return runServerlessStage("build", serverlessbuild.Run)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var serverlessDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "jzero serverless delete",
	Long:  `jzero serverless delete.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		embeded.Home = config.C.Home
		return runServerlessStage("delete", serverlessdelete.Run)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func GetCommand() *cobra.Command {
	serverlessCmd.AddCommand(serverlessBuildCmd)
	serverlessCmd.AddCommand(serverlessDeleteCmd)

	serverlessDeleteCmd.Flags().StringSliceP("plugin", "p", nil, "plugin name")

	return serverlessCmd
}

func runServerlessStage(kind string, fn func() ([]string, error)) error {
	items, err := fn()
	if config.C.Quiet {
		return err
	}

	title := console.Green(stringsTitle(kind)) + " " + console.Yellow("serverless")
	fmt.Printf("%s\n", console.BoxHeader("", title))

	for _, item := range items {
		fmt.Printf("%s\n", console.BoxItem(item))
	}

	if err != nil {
		for _, line := range console.NormalizeErrorLines(err.Error()) {
			fmt.Printf("%s\n", console.BoxDetailItem(line))
		}
		fmt.Printf("%s\n\n", console.BoxErrorFooter())
		return console.MarkRenderedError(err)
	}

	fmt.Printf("%s\n\n", console.BoxSuccessFooter())
	return nil
}

func stringsTitle(kind string) string {
	switch kind {
	case "build":
		return "Build"
	case "delete":
		return "Delete"
	default:
		return kind
	}
}
