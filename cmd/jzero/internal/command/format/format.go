/*
Copyright © 2025 jaronnie <jaron@jaronnie.com>
*/

package format

import (
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/format/formatgo"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "code format tool",
	Long:  `used to format code`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := formatgo.Run(); err != nil {
			return err
		}
		return nil
	},
	SilenceUsage: true,
}

// formatGoCmd represents the format go code command
var formatGoCmd = &cobra.Command{
	Use:   "go",
	Short: "used to format go code",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := formatgo.Run(); err != nil {
			return err
		}
		return nil
	},
	SilenceUsage: true,
}

func GetCommand() *cobra.Command {
	formatCmd.PersistentFlags().BoolP("git-change", "", true, "just format git changed files")
	formatCmd.PersistentFlags().BoolP("display-diff", "d", false, "display diffs instead of rewriting files")

	formatCmd.AddCommand(formatGoCmd)
	return formatCmd
}
