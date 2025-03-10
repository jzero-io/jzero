/*
Copyright © 2025 jaronnie <jaron@jaronnie.com>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/internal/migrate/migratedown"
	"github.com/jzero-io/jzero/internal/migrate/migrategoto"
	"github.com/jzero-io/jzero/internal/migrate/migrateup"
	"github.com/jzero-io/jzero/internal/migrate/migrateversion"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "jzero migrate",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.UsageString())
		return nil
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "migrate up",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrateup.Run(args)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "migrate down",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migratedown.Run(args)
	},
}

var migrateGotoCmd = &cobra.Command{
	Use:   "goto",
	Short: "migrate goto",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrategoto.Run(args)
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "migrate version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrateversion.Run(args)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.PersistentFlags().StringP("source", "", "file://desc/sql_migration", "migrate source")
	_ = migrateCmd.MarkFlagRequired("source")
	migrateCmd.PersistentFlags().StringP("database", "", "mysql", "migrate database")
	_ = migrateCmd.MarkFlagRequired("database")

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateGotoCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
}
