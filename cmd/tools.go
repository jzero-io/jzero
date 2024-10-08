/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/pkg/gitstatus"
)

// toolsCmd represents the tools command
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "jzero tools",
}

var toolsGitDiffCmd = &cobra.Command{
	Use:   "git-diff",
	Short: "git diff",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := os.Getwd()
		if len(args) > 0 {
			path = args[0]
		}
		var diffFiles []string
		m, _, err := gitstatus.ChangedFiles(path, "")
		if err != nil {
			os.Exit(1)
			return
		}
		diffFiles = append(diffFiles, m...)
		if len(diffFiles) > 0 {
			os.Exit(1)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(toolsCmd)
	toolsCmd.AddCommand(toolsGitDiffCmd)
}
