/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/ivm/ivmaddapi"
	"github.com/jzero-io/jzero/internal/ivm/ivmaddproto"
	"github.com/jzero-io/jzero/internal/ivm/ivminit"
)

// ivmCmd represents the interface version manage command
var ivmCmd = &cobra.Command{
	Use:   "ivm",
	Short: `Interface version manage`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var ivmInitCmd = &cobra.Command{
	Use:   "init",
	Short: `Init newer version from older version, no need to do any more`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ivminit.Init(config.C.Ivm)
	},
}

var ivmAddCmd = &cobra.Command{
	Use:   "add",
	Short: `Add example interface descriptor files`,
}

var ivmAddProtoCmd = &cobra.Command{
	Use:   "proto",
	Short: `Add a example proto`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !strings.HasPrefix(config.C.Ivm.Version, "v") {
			cobra.CheckErr(errors.New("version must has prefix v"))
		}
		if len(config.C.Ivm.Add.Proto.Services) == 0 {
			config.C.Ivm.Add.Proto.Services = []string{config.C.Ivm.Add.Proto.Name}
		}
		return ivmaddproto.AddProto(config.C)
	},
	SilenceUsage: true,
}

var ivmAddApiCmd = &cobra.Command{
	Use:   "api",
	Short: `Add a example api`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !strings.HasPrefix(config.C.Ivm.Version, "v") {
			cobra.CheckErr(errors.New("version must has prefix v"))
		}
		if config.C.Ivm.Add.Api.Group == "" {
			config.C.Ivm.Add.Api.Group = config.C.Ivm.Add.Api.Name
		}
		return ivmaddapi.AddApi(config.C)
	},
	SilenceUsage: true,
}

func init() {
	{
		rootCmd.AddCommand(ivmCmd)
		ivmCmd.PersistentFlags().StringP("version", "v", "v1", "jzero ivm version")
		ivmCmd.PersistentFlags().BoolP("split-api-types-dir", "", false, "")
	}

	{
		ivmCmd.AddCommand(ivmInitCmd)

		ivmInitCmd.Flags().StringP("style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
		ivmInitCmd.Flags().BoolP("remove-suffix", "", true, "remove suffix Handler and Logic on filename or file content")
		ivmInitCmd.Flags().BoolP("change-logic-types", "", true, "if api file or proto change, e.g. Request or Response type, change handler and logic file content types but not file")
	}

	{
		ivmCmd.AddCommand(ivmAddCmd)
	}

	{
		ivmAddCmd.AddCommand(ivmAddProtoCmd)

		ivmAddProtoCmd.Flags().StringP("name", "n", "template", "set proto file name")
		_ = ivmAddProtoCmd.MarkFlagRequired("name")

		ivmAddProtoCmd.Flags().StringSliceP("services", "", nil, "set proto services")
		ivmAddProtoCmd.Flags().StringSliceP("methods", "m", []string{"get:SayHello"}, "set proto methods")
	}

	{
		ivmAddCmd.AddCommand(ivmAddApiCmd)
		ivmAddApiCmd.Flags().StringP("name", "n", "template", "set api file name")
		_ = ivmAddApiCmd.MarkFlagRequired("name")

		ivmAddApiCmd.Flags().StringP("group", "", "", "set api file group")
		ivmAddApiCmd.Flags().StringSliceP("handlers", "", []string{"get:List", "get:Get", "post:Edit", "get:Delete"}, "set api file handlers")
	}
}
