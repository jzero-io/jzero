package add

import (
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/add/addapi"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/add/addproto"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/add/addsql"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add",
	Short:        `Used to add api/proto/sql file`,
	SilenceUsage: true,
}

var addApiCmd = &cobra.Command{
	Use:   "api",
	Short: `Add api`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return addapi.Run(args)
	},
	SilenceUsage: true,
}

var addProtoCmd = &cobra.Command{
	Use:   "proto",
	Short: `Add proto`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return addproto.Run(args)
	},
	SilenceUsage: true,
}

var addSqlCmd = &cobra.Command{
	Use:   "sql",
	Short: `Add sql`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return addsql.Run(args)
	},
	SilenceUsage: true,
}

func GetCommand() *cobra.Command {
	addCmd.PersistentFlags().StringP("output", "o", "file", "Output format. One of: file | stdout")

	addCmd.AddCommand(addApiCmd)
	addCmd.AddCommand(addProtoCmd)
	addCmd.AddCommand(addSqlCmd)
	return addCmd
}
