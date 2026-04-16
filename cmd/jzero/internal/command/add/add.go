package add

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/add/addapi"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/add/addproto"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/add/addsql"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
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
		return runAddStage("api", args, addapi.Run)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var addProtoCmd = &cobra.Command{
	Use:   "proto",
	Short: `Add proto`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAddStage("proto", args, addproto.Run)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var addSqlCmd = &cobra.Command{
	Use:   "sql",
	Short: `Add sql`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAddStage("sql", args, addsql.Run)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func GetCommand() *cobra.Command {
	addCmd.PersistentFlags().StringP("output", "o", "file", "Output format. One of: file | stdout")

	addCmd.AddCommand(addApiCmd)
	addCmd.AddCommand(addProtoCmd)
	addCmd.AddCommand(addSqlCmd)
	return addCmd
}

func runAddStage(kind string, args []string, fn func([]string) (string, error)) error {
	target, err := fn(args)
	if config.C.Add.Output != "file" || config.C.Quiet {
		return err
	}

	title := console.Green("Add") + " " + console.Yellow(kind)
	fmt.Printf("%s\n", console.BoxHeader("", title))

	if err != nil {
		if target != "" {
			fmt.Printf("%s\n", console.BoxErrorItem(target))
		}
		for _, line := range console.NormalizeErrorLines(err.Error()) {
			fmt.Printf("%s\n", console.BoxDetailItem(line))
		}
		fmt.Printf("%s\n\n", console.BoxErrorFooter())
		if config.C.Quiet {
			return err
		}
		return console.MarkRenderedError(err)
	}

	if target != "" {
		fmt.Printf("%s\n", console.BoxItem(target))
	}
	fmt.Printf("%s\n\n", console.BoxSuccessFooter())
	return nil
}
