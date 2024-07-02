/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"strings"

	"github.com/jzero-io/jzero/internal/ivm"
	"github.com/jzero-io/jzero/internal/ivm/ivmaddproto"
	"github.com/jzero-io/jzero/internal/ivm/ivminit"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	RunE:  ivminit.Init,
}

var ivmAddCmd = &cobra.Command{
	Use:   "add",
	Short: `Add example interface descriptor files`,
}

var ivmAddProtoCmd = &cobra.Command{
	Use:   "proto",
	Short: `Add a example proto`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !strings.HasPrefix(ivm.Version, "v") {
			cobra.CheckErr(errors.New("version must has prefix v"))
		}
		if len(ivmaddproto.Services) == 0 {
			ivmaddproto.Services = []string{ivmaddproto.Name}
		}
	},
	RunE:         ivmaddproto.AddProto,
	SilenceUsage: true,
}

func init() {
	{
		rootCmd.AddCommand(ivmCmd)
		ivmCmd.PersistentFlags().StringVarP(&ivm.Version, "version", "v", "v1", "jzero ivm init")
	}

	{
		ivmCmd.AddCommand(ivmInitCmd)

		ivmInitCmd.Flags().StringVarP(&ivminit.Style, "style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
		ivmInitCmd.Flags().BoolVarP(&ivminit.RemoveSuffix, "remove-suffix", "", false, "remove suffix Handler and Logic on filename or file content")
	}

	{
		ivmCmd.AddCommand(ivmAddCmd)
	}

	{
		ivmAddCmd.AddCommand(ivmAddProtoCmd)

		ivmAddProtoCmd.Flags().StringVarP(&ivmaddproto.Name, "name", "", "template", "set proto name")
		ivmAddProtoCmd.Flags().StringSliceVarP(&ivmaddproto.Services, "services", "", nil, "set proto services")
		ivmAddProtoCmd.Flags().StringSliceVarP(&ivmaddproto.Methods, "methods", "m", []string{"SayHello:get"}, "set proto methods")
	}
}
