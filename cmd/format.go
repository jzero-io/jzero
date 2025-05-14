/*
Copyright © 2025 jaronnie <jaron@jaronnie.com>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/internal/format/formatapi"
	"github.com/jzero-io/jzero/internal/format/formatgo"
	"github.com/jzero-io/jzero/internal/format/formatproto"
	"github.com/jzero-io/jzero/internal/format/formatsql"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "jzero code format tool",
	Long:  `used to format code. e.g. go/api/proto/sql`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := formatgo.Run(); err != nil {
			return err
		}

		if err := formatapi.Run(); err != nil {
			return err
		}

		if err := formatproto.Run(); err != nil {
			return err
		}

		if err := formatsql.Run(); err != nil {
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

var formatApiCmd = &cobra.Command{
	Use:   "api",
	Short: "used to format api code",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("format api called")
	},
	SilenceUsage: true,
}

var formatProtoCmd = &cobra.Command{
	Use:   "proto",
	Short: "used to format proto code",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("format proto called")
	},
	SilenceUsage: true,
}

var formatSqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "used to format sql code",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("format sql called")
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(formatCmd)
	formatCmd.PersistentFlags().BoolP("git-change", "", true, "just format git changed files")
	formatCmd.PersistentFlags().BoolP("display-diff", "d", false, "display diffs instead of rewriting files")

	formatCmd.AddCommand(formatGoCmd)
	formatCmd.AddCommand(formatApiCmd)
	formatCmd.AddCommand(formatProtoCmd)
	formatCmd.AddCommand(formatSqlCmd)
}
