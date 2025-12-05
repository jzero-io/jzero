/*
Copyright © 2025 jaronnie <jaron@jaronnie.com>
*/

package migrate

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/migrate/migratedown"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/migrate/migrategoto"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/migrate/migrateup"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/migrate/migrateversion"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate model by desc/sql_migration",
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

func GetCommand() *cobra.Command {
	migrateCmd.PersistentFlags().StringP("source", "", "file://desc/sql_migration", "migrate source")
	_ = migrateCmd.MarkFlagRequired("source")
	migrateCmd.PersistentFlags().StringP("datasource-url", "", "", "migrate datasource url")
	_ = migrateCmd.MarkFlagRequired("datasource-url")
	migrateCmd.PersistentFlags().StringP("x-migrations-table", "", "schema_migrations", "migrate table name")
	migrateCmd.PersistentFlags().BoolP("source-append-driver", "", false, "migrate source append driver")

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateGotoCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	return migrateCmd
}
