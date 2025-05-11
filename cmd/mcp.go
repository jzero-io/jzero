/*
Copyright © 2025 jaronnie <jaron@jaronnie.com>
*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/internal/mcp"
)

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "mcp server for jzero",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcp.Run(rootCmd)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
