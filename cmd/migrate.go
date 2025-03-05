/*
Copyright © 2025 jaronnie <jaron@jaronnie.com>
*/

package cmd

import (
	"github.com/jzero-io/jzero/internal/migrate"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "jzero migrate",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrate.Run()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
